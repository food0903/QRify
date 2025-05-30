pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh '''
                    docker-compose build
                '''
            }
        }
        stage('Test') {
            steps {
                sh '''
                    docker compose build backend
                    docker compose run --rm backend go test -v ./internal/tests
                '''
            }
        }
        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh '''
                    cd QRifier &&
                    git pull &&
                    docker compose down &&
                    docker compose up --build -d
                    '
                '''
            }
        }
    }
}