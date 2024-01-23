# Real Notification

Real Notification is a robust and scalable notification service designed to deliver reliable performance. This project leverages a powerful tech stack, including Apache Kafka for efficient message queuing, Postgres for persistent storage of notifications, and Gin as a framework for handling incoming requests.
![realNotification](https://github.com/SanjaySinghRajpoot/realNotification/assets/67458417/ae540187-db45-4a6a-ab44-2b7f7f769327)

## Features

### 1. Scalable Architecture

Real Notification boasts a scalable architecture, achieved through the use of Apache Kafka as a messaging queue and an MVC (Model-View-Controller) design pattern. This ensures that the system can handle varying loads of notifications efficiently.

### 2. Docker Setup

The project includes a Docker setup for both the database and the entire application. This facilitates seamless deployment and ensures consistent environments across different systems.

### 3. Fast Performance

Real Notification is engineered for fast performance, minimizing downtime even under heavy request loads. The combination of Apache Kafka for message queuing and the use of a lightweight framework like Gin ensures swift response times.

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

