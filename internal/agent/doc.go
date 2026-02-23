/*
This package provides the agent for collecting metrics and sending them to the server

When called agent methods start:

	metricAgent := agent.NewMetricsAgent(agentConfig)
	metricAgent.Start()

here 2 goroutines are started:
The first is for collection of metrics
The second is form sending metrics to server

All agent parameters are described in config package
*/
package agent
