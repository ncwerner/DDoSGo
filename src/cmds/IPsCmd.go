package cmds

import "subsytems"
import "visitors"

type IPsCmd struct{
     agent *subsystems.Enforcer
     agent_ip string
     agent_port string
     OffendingIPs []string
}

func NewIPsCmd(aIP string, aPort string, ips []string) IPsCmd{
     return IPsCmd{
     	    agent: visitors.AgentRegVisitor.Agent
	    agent_ip: aIP,
	    agent_port: aPort,
	    OffendingIPs: ips}

//Need a way to get an enforcer instance into the New IP cmd
//Maybe pass that in when it is created otherwise we would have to create the instance from within this 'package'