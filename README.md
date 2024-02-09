# Real Notification

Real Notification is a robust and scalable notification service designed to deliver reliable performance. This project leverages a powerful tech stack, including Apache Kafka for efficient message queuing, Postgres for persistent storage of notifications, and Gin as a framework for handling incoming requests.

![systemdesignDiagram](https://github.com/SanjaySinghRajpoot/realNotification/assets/67458417/5d3d5761-e05d-465c-b8b8-b7456d92003a)



## Features

- Utilize Apache Kafka as a robust messaging queue for efficient handling of heavy messages.
- Employ Docker Compose for a comprehensive project setup.
- Ensure swift response times with the lightweight Gin framework.
- Implement an IP-based rate limiter to prevent spam requests.
- Organize API endpoints using an MVC architecture.
- Embrace a microservices architecture for all services.
- Conducted load testing using k6.
- Implement a CRON job to identify and resend failed notifications.
- Integrated Redis Cache for checking repeated Notifications
- Used Database sharding for saving the notifications

- Used my [LoadBalancer](https://github.com/SanjaySinghRajpoot/loadbalancer) Project for handling traffic 
## Getting Started

To get started with Real Notification, follow these steps:

1. **Clone the Repository**: `git clone [repository_url]`

2. **Build Docker Containers**: Navigate to the project directory and run `docker-compose up --build` to set up the Docker containers for the database and the application.

3. **Configure the System**: Adjust configuration files as needed, specifying connection details for Apache Kafka and Postgres.

4. **Run the Application**: Start the application using `go run main.go` or as per your preferred deployment method.


## Contributing

We welcome contributions in the form of bug reports, feature requests, or pull requests.

## License

Real Notification is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute the code as per the terms of the license.

