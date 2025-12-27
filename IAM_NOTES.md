# IAM Notes for AWS Admin

This app calls AWS APIs from the backend. Attach only the services you need.
Replace `*` with tighter ARNs/regions where possible.

## EC2 + VPC (includes subnets, security groups, key pairs)

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "EC2FullManagement",
      "Effect": "Allow",
      "Action": [
        "ec2:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## S3

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "S3FullManagement",
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## ECR

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ECRFullManagement",
      "Effect": "Allow",
      "Action": [
        "ecr:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## ECS

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ECSFullManagement",
      "Effect": "Allow",
      "Action": [
        "ecs:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Lambda

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "LambdaFullManagement",
      "Effect": "Allow",
      "Action": [
        "lambda:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## CloudFormation

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "CloudFormationFullManagement",
      "Effect": "Allow",
      "Action": [
        "cloudformation:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## CloudWatch Logs

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "CloudWatchLogsFullManagement",
      "Effect": "Allow",
      "Action": [
        "logs:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Route 53

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Route53FullManagement",
      "Effect": "Allow",
      "Action": [
        "route53:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## ACM (Certificate Manager)

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AcmFullManagement",
      "Effect": "Allow",
      "Action": [
        "acm:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Secrets Manager

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "SecretsManagerFullManagement",
      "Effect": "Allow",
      "Action": [
        "secretsmanager:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Systems Manager (SSM)

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "SSMFullManagement",
      "Effect": "Allow",
      "Action": [
        "ssm:*",
        "ssmmessages:*",
        "ec2messages:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## ELBv2 (Load Balancer)

Full management (broad):

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ELBv2FullManagement",
      "Effect": "Allow",
      "Action": [
        "elasticloadbalancing:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Combined Policy (All App Services, Regions: US/EU/Asia)

This policy allows all services used by this app, restricted to regions in
US, EU, and Asia-Pacific (including China and GovCloud exclusions by default).
Route 53 is global and not region-scoped, so it remains allowed.

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AppServicesAllRegions",
      "Effect": "Allow",
      "Action": [
        "acm:*",
        "cloudformation:*",
        "ec2:*",
        "ecr:*",
        "ecs:*",
        "elasticloadbalancing:*",
        "lambda:*",
        "logs:*",
        "s3:*",
        "secretsmanager:*",
        "ssm:*",
        "ssmmessages:*",
        "ec2messages:*"
      ],
      "Resource": "*",
      "Condition": {
        "StringLike": {
          "aws:RequestedRegion": [
            "us-*",
            "eu-*",
            "ap-*"
          ]
        }
      }
    },
    {
      "Sid": "Route53Global",
      "Effect": "Allow",
      "Action": [
        "route53:*"
      ],
      "Resource": "*"
    }
  ]
}
```

## Notes

- Start with only the services you actually use.
- Scope by region and resource ARNs whenever possible.
- Do not use root credentials.

## Client Role Trust Policy (template)

Replace `YOUR_BACKEND_ACCOUNT_ID` and `EXTERNAL_ID_VALUE`.

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::YOUR_BACKEND_ACCOUNT_ID:root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "EXTERNAL_ID_VALUE"
        }
      }
    }
  ]
}
```
