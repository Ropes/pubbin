# pubbing Google PubSub debugging CLI


## Example Commands

* Publish a single hello world debug message to topic:
  * `./pubbing pub --project=<project> --topic=<topic>

* Subscribe to messages and give verbose output:
  * `./pubbing sub --project=<project> --topic=<topic> --sub=<subname> --num=1000  --log=debug`
  * Log level `debug` will write the messages to stdout
  * Log level `info` is meant to test performance and will print subscription speed metrics

* Publish messages from [LegendaryGopher](https://github.com/schmichael/legendarygopher) Figures API to a topic
  * `./pubbing gopherpump  --project=<project> --topic=<topic>  --num=20000  --batch=1000` 

