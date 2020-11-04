# publisher
The Lunara Automatic Publisher

## features
- Abstract interface StepOperator: includes ftp, git, and svn operator, and they
all implement the github.com/nevercase/publisher/pkg/intefaces.StepOperator
- Scheduler: the center of the whole system, which supplies a series of apis about demonstrating the
 dashboard and controlling all the runners
- Runner: contains multiple k-v values which were used to control the Runner to take actions.