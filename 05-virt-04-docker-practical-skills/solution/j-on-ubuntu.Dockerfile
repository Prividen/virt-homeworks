FROM ubuntu:latest

RUN apt-get update && apt-get -y install gnupg wget && apt-get -y install ca-certificates && \
	wget -q -O - https://pkg.jenkins.io/debian-stable/jenkins.io.key | apt-key add - && \
	echo "deb https://pkg.jenkins.io/debian-stable binary/" >> /etc/apt/sources.list && \
	apt-get update && apt-get -y install openjdk-11-jre jenkins && \
	apt-get clean

EXPOSE 8080

USER jenkins
ENTRYPOINT ["java"]
CMD ["-Djava.awt.headless=true", "-jar", "/usr/share/jenkins/jenkins.war"]
