pipeline {
    agent any
    stages {
        stage('Build & Test') {
            when {
                not {
                    branch 'main'
                }
            }
            steps {
                dir('/var/lib/jenkins/QRify') {
                    sh '''
                        docker compose build
                        docker compose run --rm backend go test -v ./internal/tests
                    '''
                }
            }
        }
        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                dir('/var/lib/jenkins/QRify') {
                    sh '''
                        git pull
                        docker compose down
                        docker compose up --build -d
                    '''
                }
            }
        }
    }
}