
# Go Mircoservices - Application

Hello, below in the README file is the current state of the application based on microservices architecture. Currently the application has 3 services:

### Front-end Service

The appearance of the website consists of a simple layout which is divided into three sections. The first is the Test Broker, which contains a button that connects to the port on which the Broker service is served and the response it gives us. The next two sections run automatically when the button is pressed. The first is what was sent to the service in question, and the next is what data was captured ( raw )

![App Screenshot](https://github.com/kleczynski/GoMicroservices/blob/master/images/layout.png)

### Broker Service 

In this service we have routers that are responsible for communication via the REST API with other services, at the moment the service broker's communication with the front-end is tested. The service broker will distribute the communication with the front-end to the individual services that are planned. 

### Authentication Service

The last service is currently responsible for authorising future users. The current image is integrated into a Postgresql database, the image of which, as well as the other services, can be found in docker compose. 

# Legacy

- Logger Service
- Listener Service with RabbitMQ
- Communication between servives using gRPC
- Deploy Distributed app on K8's 

# Building Application

To build or start application we have to be in ```project``` dir. There we have ```Makefile``` We can use command such as. 
For the first time ```make up_build``` it will build 3 images with current services that are available inside project. To stop services we can use ```make down```
