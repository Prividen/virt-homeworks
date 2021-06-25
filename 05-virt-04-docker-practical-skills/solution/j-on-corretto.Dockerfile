FROM amazoncorretto:11
ARG J_LINK=https://get.jenkins.io/war-stable/2.289.1/jenkins.war
RUN curl -s $(curl -sI $J_LINK |grep -m 1 link: |cut -f2 -d'<' |cut -f1 -d'>') > /jenkins.war && \
	mkdir /var/lib/jenkins
ENV JENKINS_HOME=/var/lib/jenkins
EXPOSE 8080
ENTRYPOINT ["/usr/bin/java"]
CMD ["-jar", "/jenkins.war"]
