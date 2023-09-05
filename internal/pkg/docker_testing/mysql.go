package docker_testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"strings"
	"time"
)

func RunMySQL() (string, func()) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	resp, err := c.ContainerCreate(
		ctx,
		&container.Config{
			Image: "mysql:8.0.29",
			ExposedPorts: nat.PortSet{
				"3306/tcp": {},
			},
			Env: []string{"MYSQL_ROOT_PASSWORD=my-secret-pw"},
		},
		&container.HostConfig{
			Binds: []string{
				"/Users/fs/Desktop/yellowbook/script/mysql/:/docker-entrypoint-initdb.d/",
			},
			PortBindings: nat.PortMap{
				"3306/tcp": []nat.PortBinding{
					{
						HostIP:   "127.0.0.1",
						HostPort: "0",
					},
				},
			},
		},
		nil,
		nil,
		"unit-test-mysql",
	)

	if err != nil {
		panic(err)
	}

	containerId := resp.ID

	err = c.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	inspRes, err := c.ContainerInspect(ctx, containerId)
	if err != nil {
		panic(err)
	}

	hostPort := inspRes.NetworkSettings.Ports["3306/tcp"][0]

	for {
		err := checkMySQLService(ctx, c, containerId)

		if err == nil {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Sprintf("root:my-secret-pw@tcp(%s:%s)/yellowbook", hostPort.HostIP, hostPort.HostPort), func() {
		err := c.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}
}

func checkMySQLService(ctx context.Context, cli *client.Client, containerId string) error {
	inspect, err := cli.ContainerInspect(ctx, containerId)
	if err != nil {
		return err
	}

	if !inspect.State.Running {
		return errors.New("容器已退出")
	}

	logOutput, err := getContainerLogs(ctx, cli, containerId)
	if err != nil {
		return err
	}

	if !strings.Contains(logOutput, "/usr/sbin/mysqld: ready for connections. Version: '8.0.29'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server - GPL") {
		return errors.New("MySQL 未启动完成")
	}

	return nil
}

func getContainerLogs(ctx context.Context, cli *client.Client, containerId string) (string, error) {
	logReader, err := cli.ContainerLogs(ctx, containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "all",
	})

	if err != nil {
		return "", err
	}
	defer logReader.Close()

	logOutput := ""
	logBuffer := make([]byte, 4096)

	for {
		n, err := logReader.Read(logBuffer)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		logOutput += string(logBuffer[:n])
	}

	return logOutput, nil
}
