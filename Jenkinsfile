pipeline {
    agent any
    stages {
        stage('Test') {
            steps {
                sh '''
                    cd backend
                    go mod download
                    go test -v ./internal/tests
                '''
            }
        }
    }
}