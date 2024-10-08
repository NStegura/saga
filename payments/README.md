kafka-ui 
http://localhost:8080/

консоль
kafkacat -C -b localhost:9095 -t payment
kafkacat -C -b localhost:9095 -t inventory

консоль redis  
telnet localhost 6379