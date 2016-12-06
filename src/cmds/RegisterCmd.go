package cmds

import "subsystems"
import "visitors"

type RegisterCmd struct {
	agentRegistry *subsystems.AgentRegistry
	agent_ip      string
	handler_ip    string
	handler_port  string
	agent_port    string
	traceroute    []string
	addRemoveFlag bool
}

func NewRegisterCmd(aIP string, hIP string, hPort string, trace []string, aPort string, flag bool) RegisterCmd {
	return RegisterCmd{
		agentRegistry: visitors.AgentRegVisitor.AgentReg,
		agent_ip:      aIP,
		handler_ip:    hIP,
		handler_port:  hPort,
		agent_port:    aPort,
		traceroute:    trace,
		addRemoveFlag: flag}
}

func (regcmd RegisterCmd) ExecCmd() {
	//fmt.Printf("RegCmd internals: %s %s %s %t\n", regcmd.agent_ip, regcmd.handler_ip, regcmd.handler_port, regcmd.addRemoveFlag)

	//false = add agent to registry
	//true = remove agent from registry

	//Construct an agent record from given information
	record := subsystems.NewAgentRecord(regcmd.agent_ip, regcmd.handler_ip, regcmd.handler_port, regcmd.agent_port, regcmd.traceroute)

	if !regcmd.addRemoveFlag {
		regcmd.agentRegistry.AddAgent(*record)
	} else {
		regcmd.agentRegistry.RemoveAgent(regcmd.agent_ip, regcmd.agent_port)
	}
}
