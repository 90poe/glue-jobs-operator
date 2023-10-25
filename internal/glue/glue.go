package glue

import (
	"context"
	"fmt"
	"strings"

	awsv1alpha1 "github.com/90poe/glue-jobs-operator/api/v1alpha1"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsglue "github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
)

const (
	glueETL = "glueetl"
)

var (
	// jobOwner is map which we will add to Tags so we know,
	// that this GlueJob is owned by this operator
	jobOwnedByOperator = map[string]string{
		"glue-jobs-operator": "true",
	}
)

type Job struct {
	ctx       context.Context
	job       awsv1alpha1.GlueJobSpec
	exists    bool
	awsClient *awsglue.Client
}

// NewJob will return a new Job struct
func NewJob(ctx context.Context, job awsv1alpha1.GlueJobSpec) (*Job, error) {
	gJob := &Job{
		ctx:    ctx,
		job:    job,
		exists: false,
	}
	// Load the Shared AWS Configuration from the Shared config file or IAM Roles
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Create an Amazon Glue service client
	gJob.awsClient = awsglue.NewFromConfig(cfg)

	// check that GlueJob exists on AWS
	gJob.exists, err = gJob.checkJobExistsOnAWS()
	if err != nil {
		return nil, fmt.Errorf("failed to check if GlueJob %s exists on AWS: %w", job.Name, err)
	}

	return gJob, nil
}

// JobExists will return true if GlueJob exists on AWS
func (g *Job) JobExists() bool {
	return g.exists
}

// CreateJob will create Glue Job
func (g *Job) CreateJob() error {
	command := &types.JobCommand{
		Name:           aws.String(g.job.Command.Name),
		PythonVersion:  aws.String(fmt.Sprintf("%v", g.job.Command.PythonVersion)),
		ScriptLocation: aws.String(g.job.Command.ScriptLocation),
	}
	if strings.ToLower(g.job.Command.Name) != glueETL {
		command.Runtime = aws.String(g.job.Command.Runtime)
	}
	input := &awsglue.CreateJobInput{
		Name:            aws.String(g.job.Name),
		Command:         command,
		Role:            aws.String(g.job.Role),
		Timeout:         aws.Int32(g.job.TimeoutInMinutes),
		GlueVersion:     aws.String(g.job.GlueVersion),
		NumberOfWorkers: aws.Int32(g.job.NumberOfWorkers),
		WorkerType:      types.WorkerType(g.job.WorkerType),
		ExecutionClass:  types.ExecutionClass(g.job.ExecutionClass),
		ExecutionProperty: &types.ExecutionProperty{
			MaxConcurrentRuns: g.job.ExecutionProperty.MaxConcurrentRuns,
		},
		MaxRetries:       g.job.MaxRetries,
		DefaultArguments: g.job.DefaultArguments,
		Tags:             jobOwnedByOperator,
	}
	// create job
	_, err := g.awsClient.CreateJob(g.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create Glue Job %s: %w", g.job.Name, err)
	}
	return nil
}

// UpdateJob will update Glue Job
func (g *Job) UpateJob() error {
	command := &types.JobCommand{
		Name:           aws.String(g.job.Command.Name),
		PythonVersion:  aws.String(fmt.Sprintf("%v", g.job.Command.PythonVersion)),
		ScriptLocation: aws.String(g.job.Command.ScriptLocation),
	}
	if strings.ToLower(g.job.Command.Name) != glueETL {
		command.Runtime = aws.String(g.job.Command.Runtime)
	}
	input := &awsglue.UpdateJobInput{
		JobName: aws.String(g.job.Name),
		JobUpdate: &types.JobUpdate{
			Command:         command,
			Role:            aws.String(g.job.Role),
			Timeout:         aws.Int32(g.job.TimeoutInMinutes),
			GlueVersion:     aws.String(g.job.GlueVersion),
			NumberOfWorkers: aws.Int32(g.job.NumberOfWorkers),
			WorkerType:      types.WorkerType(g.job.WorkerType),
			ExecutionClass:  types.ExecutionClass(g.job.ExecutionClass),
			ExecutionProperty: &types.ExecutionProperty{
				MaxConcurrentRuns: g.job.ExecutionProperty.MaxConcurrentRuns,
			},
			MaxRetries:       g.job.MaxRetries,
			DefaultArguments: g.job.DefaultArguments,
		},
	}
	// update job
	_, err := g.awsClient.UpdateJob(g.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update Glue Job %s: %w", g.job.Name, err)
	}
	return nil
}

// DeleteJob will delete Glue Job
func (g *Job) DeleteJob() error {
	if !g.exists {
		return nil
	}
	// delete job
	_, err := g.awsClient.DeleteJob(g.ctx, &awsglue.DeleteJobInput{
		JobName: aws.String(g.job.Name),
	})
	if err != nil {
		return fmt.Errorf("failed to delete Glue Job %s: %w", g.job.Name, err)
	}
	return nil
}

func (g *Job) checkJobExistsOnAWS() (bool, error) {
	// get all jobs paginator
	jobsPaginator := awsglue.NewListJobsPaginator(g.awsClient, &awsglue.ListJobsInput{
		MaxResults: aws.Int32(100),
		Tags:       jobOwnedByOperator,
	})
	// Iterate through all jobs in paginator
	for {
		if !jobsPaginator.HasMorePages() {
			// no more pages, break
			break
		}
		jobsOut, err := jobsPaginator.NextPage(g.ctx)
		if err != nil {
			return false, fmt.Errorf("failed to get next page from Jobs paginator: %w", err)
		}
		// we have 2 criteria to match a job:
		// 1. job name
		// 2. job tags
		for _, jobName := range jobsOut.JobNames {
			if jobName == g.job.Name {
				return true, nil
			}
		}
	}
	return false, nil
}
