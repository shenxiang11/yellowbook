package cloopen

// 真的会发短信，所以注释掉吧

//func TestCloopen(t *testing.T) {
//	cfg := cloopen.DefaultConfig().
//		WithAPIAccount("8aaf07087fe90a32017ff389d6ac01bb").
//		WithAPIToken("a1c23065a7d847c384d719ad240f6384")
//
//	client := cloopen.NewJsonClient(cfg)
//
//	s := NewService(client)
//
//	err := s.Send(context.TODO(), "1", []sms.NamedArg{
//		{
//			Name: "1",
//			Val:  "1234",
//		},
//		{
//			Name: "2",
//			Val:  "25444444",
//		},
//	}, "18616154465")
//
//	if err != nil {
//		t.Errorf("Error %v", err)
//	}
//}
