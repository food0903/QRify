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
                        docker image prune -a -f
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
                        
                        docker compose up --build -d --scale backend=2 --scale frontend=2
                        
                        docker image prune -a -f
                    '''
                }
            }
        }
    }
}