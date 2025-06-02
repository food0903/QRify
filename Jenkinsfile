pipeline {
    agent any
    stages {
        stage('Build & Test') {
            steps {
                sh '''
                    cd QRify
                    docker compose build
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
                    cd QRify
                    git pull
                    docker compose down
                    docker compose up --build -d
                '''
            }
        }
    }
}