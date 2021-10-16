//resource "aws_ecr_repository" "scheduler_repo" {
//  name = "scheduler"
//}
//
//resource "aws_ecr_repository" "videoservice_repo" {
//  name = "videoservice"
//}
//
//resource "aws_ecr_repository" "userservice_repo" {
//  name = "userservice"
//}
//
//resource "aws_ecr_repository" "frontend_repo" {
//  name = "frontend"
//}

resource "aws_ecr_repository" "react_frontend" {
  name = "react_frontend"
}