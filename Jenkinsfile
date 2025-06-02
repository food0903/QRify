pipeline {
    agent any
    stages {
        stage('Build & Test') {
            steps {
                dir('/root/QRify') {
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
                dir('/root/QRify') {
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