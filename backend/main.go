package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bertrandmartel/aws-admin/backend/api"
	awsMiddleware "github.com/bertrandmartel/aws-admin/backend/middleware"
	"github.com/bertrandmartel/aws-admin/backend/model"
	"github.com/bertrandmartel/aws-admin/backend/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

var outputTmpDir = ""

const rootTmpDir = "tmp"

func main() {
	config, err := model.ParseConfig("config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	if os.Getenv("AWS_REGION") != "" {
		config.DefaultRegion = os.Getenv("AWS_REGION")
	}
	if _, err := os.Stat(rootTmpDir); os.IsNotExist(err) {
		os.Mkdir(rootTmpDir, os.FileMode(0755))
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	outputTmpDir, err = ioutil.TempDir(rootTmpDir, "aws-admin")
	if err != nil {
		log.Fatal(err)
		return
	}
	config.OutputTmpDir = outputTmpDir

	log.Printf("[aws-admin] version %v\n", config.Version)
	log.Printf("[aws-admin] server path %v:%v\n", config.ServerPath, config.Port)

	dbURL := os.Getenv("DB_URL")
	db, err := store.Connect(dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	migrationsPath := filepath.Join("backend", "migrations")
	if err := db.EnsureMigrations(migrationsPath); err != nil {
		migrationsPath = filepath.Join("migrations")
		if err := db.EnsureMigrations(migrationsPath); err != nil {
			log.Fatal(err)
			return
		}
	}

	e := echo.New()
	/*
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		}))
	*/
	UseCommonMiddleware(e)
	e.Use(awsMiddleware.AttachConfig(config))
	e.Use(awsMiddleware.AttachDB(db))
	routes(e, config)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
}

func routes(e *echo.Echo, cfg *model.Config) {
	e.POST("/auth/signup", api.Signup)
	e.GET("/auth/verify", api.VerifyEmail)
	e.POST("/auth/login", api.Login)
<<<<<<< HEAD
	e.GET("/auth/google", api.GoogleLogin)
	e.GET("/auth/google/callback", api.GoogleCallback)
	e.POST("/auth/invite/accept", api.AcceptInvitePublic)
	e.POST("/auth/forgot", api.ForgotPassword)
	e.POST("/auth/reset", api.ResetPassword)
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
	e.POST("/auth/refresh", api.Refresh)
	e.POST("/auth/logout", api.Logout)
	e.GET("/auth/me", api.Me, awsMiddleware.RequireAuth)

	e.POST("/configuration", api.SetConfiguration, awsMiddleware.RequireAuth)
	e.POST("/session_reset", func(c echo.Context) error {
		return c.NoContent(200)
	}, awsMiddleware.RequireAuth)

	e.GET("/tenants", api.ListTenants, awsMiddleware.RequireAuth)
	e.POST("/tenants", api.CreateTenant, awsMiddleware.RequireAuth)
	e.POST("/tenants/:id/switch", api.SwitchTenant, awsMiddleware.RequireAuth)
	e.POST("/tenants/:id/invite", api.InviteTenant, awsMiddleware.RequireAuth)
	e.POST("/tenants/accept", api.AcceptInvite, awsMiddleware.RequireAuth)

	//s3
	e.GET("/service/s3", api.GetBucketList, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.PUT("/service/s3", api.CreateBucket, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/s3", api.DeleteBucket, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/s3/objects", api.ListObjects, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.PUT("/service/s3/objects", api.CreateObject, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/s3/objects", api.DeleteObject, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/s3/objects/download", api.DownloadObject, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//elbv2
	e.GET("/service/elbv2/lb", api.GetLbList, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/elbv2/listeners", api.GetListeners, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/elbv2/tg", api.GetTargetGroups, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//acm
	e.GET("/service/acm/list", api.GetAcmList, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.PUT("/service/acm/import", api.ImportCertificates, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/acm", api.DeleteCert, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//ec2
	e.GET("/service/ec2/instances", api.GetEC2Instances, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/ec2/keypair", api.GetEC2Keypair, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.PUT("/service/ec2/keypair", api.CreateEC2Keypair, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/ec2/keypair", api.DeleteEC2Keypair, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.POST("/service/ec2/instance", api.SetEC2Status, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/ec2/vpc", api.GetVPC, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/ec2/subnets", api.GetSubnets, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/ec2/sg", api.GetSecurityGroups, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//secretsmanager
	e.GET("/service/secretsmanager", api.GetSecrets, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/secretsmanager/value", api.GetSecretValue, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.POST("/service/secretsmanager", api.CreateSecret, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/secretsmanager", api.DeleteSecret, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//lambda
	e.GET("/service/lambda", api.GetLambdas, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/lambda/function", api.GetLambdaFunction, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//cloudwatchlogs
	e.GET("/service/cloudwatch/loggroup", api.GetLogGroups, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/cloudwatch/logstreams", api.GetLogStreams, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/cloudwatch/logevents", api.GetLogEvents, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//cloudformation
	e.GET("/service/cloudformation/stacks", api.GetStacks, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.PUT("/service/cloudformation/stack", api.CreateStack, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/cloudformation/stack", api.DeleteStack, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.POST("/service/cloudformation/stack", api.UpdateStack, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/cloudformation/template", api.GetTemplate, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.POST("/service/cloudformation/termination_protection", api.UpdateTerminationProtection, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/cloudformation/description", api.GetStackDescription, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/cloudformation/events", api.GetStackEvent, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//ssm
	e.GET("/service/ssm/instances", api.ListSSM, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.POST("/service/ssm/session", api.StartSession, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//ecr
	e.GET("/service/ecr/repositories", api.GetRepositories, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/ecr/images", api.GetEcrImages, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.PUT("/service/ecr", api.CreateEcrRepo, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/ecr", api.DeleteEcrRepo, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//ecs
	e.GET("/service/ecs/clusters", api.ListECSClusters, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/ecs/services", api.ListECSServices, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/ecs/tasks", api.ListECSTasks, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.POST("/service/ecs/update", api.UpdateECSService, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	//route53
	e.GET("/service/route53/hostedzones", api.GetHostedZones, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.GET("/service/route53/recordsets", api.GetRecordSets, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.POST("/service/route53/recordsets", api.CreateRecordSet, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)
	e.DELETE("/service/route53/recordsets", api.DeleteRecordSet, awsMiddleware.RequireAuth, awsMiddleware.RequireAWSAccount)

	e.GET("/aws/accounts", api.ListAWSAccounts, awsMiddleware.RequireAuth)
	e.POST("/aws/accounts", api.CreateAWSAccount, awsMiddleware.RequireAuth)
	e.PUT("/aws/accounts/:id", api.UpdateAWSAccount, awsMiddleware.RequireAuth)
	e.DELETE("/aws/accounts/:id", api.DeleteAWSAccount, awsMiddleware.RequireAuth)

	// admin (super admin)
	e.GET("/admin/users", api.ListUsers, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.POST("/admin/users", api.CreateUser, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.PUT("/admin/users/:id", api.UpdateUser, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.DELETE("/admin/users/:id", api.DeleteUser, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.POST("/admin/users/:id/reset_password", api.ResetUserPassword, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.POST("/admin/users/:id/verify_email", api.VerifyUserEmail, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)

	e.GET("/admin/tenants", api.ListTenantsAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.POST("/admin/tenants", api.CreateTenantAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.PUT("/admin/tenants/:id", api.UpdateTenantAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.DELETE("/admin/tenants/:id", api.DeleteTenantAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.GET("/admin/tenants/:id/users", api.ListTenantMembersAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.POST("/admin/tenants/:id/users", api.AddTenantMemberAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.PUT("/admin/tenants/:id/users/:userId", api.UpdateTenantMemberAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.DELETE("/admin/tenants/:id/users/:userId", api.RemoveTenantMemberAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)

	e.GET("/admin/tenants/:id/aws/accounts", api.ListAWSAccountsForTenantAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.POST("/admin/tenants/:id/aws/accounts", api.CreateAWSAccountAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.PUT("/admin/tenants/:id/aws/accounts/:accountId", api.UpdateAWSAccountAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.DELETE("/admin/tenants/:id/aws/accounts/:accountId", api.DeleteAWSAccountAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)
	e.GET("/admin/aws/accounts", api.ListAWSAccountsAdmin, awsMiddleware.RequireAuth, awsMiddleware.RequireSuperAdmin)

	// admin (tenant admin)
	e.GET("/admin/tenant/users", api.ListTenantMembersSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)
	e.POST("/admin/tenant/users", api.AddTenantMemberSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)
	e.PUT("/admin/tenant/users/:userId", api.UpdateTenantMemberSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)
	e.DELETE("/admin/tenant/users/:userId", api.RemoveTenantMemberSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)
	e.GET("/admin/tenant/aws/accounts", api.ListAWSAccountsSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)
	e.POST("/admin/tenant/aws/accounts", api.CreateAWSAccountSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)
	e.PUT("/admin/tenant/aws/accounts/:accountId", api.UpdateAWSAccountSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)
	e.DELETE("/admin/tenant/aws/accounts/:accountId", api.DeleteAWSAccountSelf, awsMiddleware.RequireAuth, awsMiddleware.RequireTenantAdmin)

	e.GET("/admin/audit", api.ListAuditLogs, awsMiddleware.RequireAuth)

	if os.Getenv("APP_ENV") == "development" {
		devURL, err := url.Parse("http://localhost:8080")
		if err != nil {
			e.Logger.Fatal(err)
		}
		targets := []*middleware.ProxyTarget{
			{
				URL: devURL,
			},
		}
		e.Group("/js/app.js", middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
		e.Group("/*", middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	} else {
		e.Group("/js/app.js", middleware.StaticWithConfig(middleware.StaticConfig{
			Root:   "dist",
			Browse: false,
<<<<<<< HEAD
			HTML5:  true,
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
		}))
		e.Group("/*", middleware.StaticWithConfig(middleware.StaticConfig{
			Root:   "dist",
			Browse: false,
<<<<<<< HEAD
			HTML5:  true,
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
		}))
	}
}

//middleware for validation
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}


func UseCommonMiddleware(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} - - ${time_rfc3339_nano} \"${method} ${uri} ${protocol}\" ${status} ${bytes_out} \"${referer}\" \"${user_agent}\"\n",
	}))
	e.Use(middleware.Recover())
}
