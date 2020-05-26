resource "aws_iam_user" "iam_otoman" {
  name = "otoman"
  path = "/users/"
}

// TODO: make this into a group, attach users to
resource "aws_iam_user_policy" "admin_policy" {
  name = "administrator-user-policy"
  user = "${aws_iam_user.iam_otoman.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}