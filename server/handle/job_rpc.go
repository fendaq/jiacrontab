package handle

import (
	"jiacrontab/libs"
	"jiacrontab/libs/proto"
	"jiacrontab/server/conf"
	"jiacrontab/server/model"
	"log"
)

type Logic struct{}

func (l *Logic) Register(args proto.ClientConf, reply *proto.MailArgs) error {
	defer libs.MRecover()

	modelInstance := model.NewModel()
	*reply = proto.MailArgs{
		Host: conf.ConfigArgs.MailHost,
		User: conf.ConfigArgs.MailUser,
		Pass: conf.ConfigArgs.MailPass,
		Port: conf.ConfigArgs.MailPort,
	}

	modelInstance.InnerStore().Wrap(func(s *model.Store) {
		s.RpcClientList[args.Addr] = args
	}).Sync()

	log.Println("register client", args)
	return nil
}

func (l *Logic) Depends(args []proto.MScript, reply *bool) error {
	modelInstance := model.NewModel()
	log.Printf("Callee Logic.Depend taskId %s", args[0].TaskId)
	*reply = true
	for _, v := range args {
		if err := modelInstance.RpcCall(v.Dest, "Task.ExecDepend", v, &reply); err != nil {
			*reply = false
			return err
		}
	}

	return nil
}

func (l *Logic) DependDone(args proto.MScript, reply *bool) error {
	modelInstance := model.NewModel()
	log.Printf("Callee Logic.DependDone task %s", args.Name)
	*reply = true
	if err := modelInstance.RpcCall(args.Dest, "Task.ResolvedDepends", args, &reply); err != nil {
		*reply = false
		return err
	}

	return nil
}
