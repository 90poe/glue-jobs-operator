// Code generated by smithy-go-codegen DO NOT EDIT.

package glue

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates an Glue machine learning transform. This operation creates the
// transform and all the necessary parameters to train it. Call this operation as
// the first step in the process of using a machine learning transform (such as the
// FindMatches transform) for deduplicating data. You can provide an optional
// Description , in addition to the parameters that you want to use for your
// algorithm. You must also specify certain parameters for the tasks that Glue runs
// on your behalf as part of learning from your data and creating a high-quality
// machine learning transform. These parameters include Role , and optionally,
// AllocatedCapacity , Timeout , and MaxRetries . For more information, see Jobs (https://docs.aws.amazon.com/glue/latest/dg/aws-glue-api-jobs-job.html)
// .
func (c *Client) CreateMLTransform(ctx context.Context, params *CreateMLTransformInput, optFns ...func(*Options)) (*CreateMLTransformOutput, error) {
	if params == nil {
		params = &CreateMLTransformInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CreateMLTransform", params, optFns, c.addOperationCreateMLTransformMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CreateMLTransformOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CreateMLTransformInput struct {

	// A list of Glue table definitions used by the transform.
	//
	// This member is required.
	InputRecordTables []types.GlueTable

	// The unique name that you give the transform when you create it.
	//
	// This member is required.
	Name *string

	// The algorithmic parameters that are specific to the transform type used.
	// Conditionally dependent on the transform type.
	//
	// This member is required.
	Parameters *types.TransformParameters

	// The name or Amazon Resource Name (ARN) of the IAM role with the required
	// permissions. The required permissions include both Glue service role permissions
	// to Glue resources, and Amazon S3 permissions required by the transform.
	//   - This role needs Glue service role permissions to allow access to resources
	//   in Glue. See Attach a Policy to IAM Users That Access Glue (https://docs.aws.amazon.com/glue/latest/dg/attach-policy-iam-user.html)
	//   .
	//   - This role needs permission to your Amazon Simple Storage Service (Amazon
	//   S3) sources, targets, temporary directory, scripts, and any libraries used by
	//   the task run for this transform.
	//
	// This member is required.
	Role *string

	// A description of the machine learning transform that is being defined. The
	// default is an empty string.
	Description *string

	// This value determines which version of Glue this machine learning transform is
	// compatible with. Glue 1.0 is recommended for most customers. If the value is not
	// set, the Glue compatibility defaults to Glue 0.9. For more information, see
	// Glue Versions (https://docs.aws.amazon.com/glue/latest/dg/release-notes.html#release-notes-versions)
	// in the developer guide.
	GlueVersion *string

	// The number of Glue data processing units (DPUs) that are allocated to task runs
	// for this transform. You can allocate from 2 to 100 DPUs; the default is 10. A
	// DPU is a relative measure of processing power that consists of 4 vCPUs of
	// compute capacity and 16 GB of memory. For more information, see the Glue
	// pricing page (https://aws.amazon.com/glue/pricing/) . MaxCapacity is a mutually
	// exclusive option with NumberOfWorkers and WorkerType .
	//   - If either NumberOfWorkers or WorkerType is set, then MaxCapacity cannot be
	//   set.
	//   - If MaxCapacity is set then neither NumberOfWorkers or WorkerType can be set.
	//   - If WorkerType is set, then NumberOfWorkers is required (and vice versa).
	//   - MaxCapacity and NumberOfWorkers must both be at least 1.
	// When the WorkerType field is set to a value other than Standard , the
	// MaxCapacity field is set automatically and becomes read-only. When the
	// WorkerType field is set to a value other than Standard , the MaxCapacity field
	// is set automatically and becomes read-only.
	MaxCapacity *float64

	// The maximum number of times to retry a task for this transform after a task run
	// fails.
	MaxRetries *int32

	// The number of workers of a defined workerType that are allocated when this task
	// runs. If WorkerType is set, then NumberOfWorkers is required (and vice versa).
	NumberOfWorkers *int32

	// The tags to use with this machine learning transform. You may use tags to limit
	// access to the machine learning transform. For more information about tags in
	// Glue, see Amazon Web Services Tags in Glue (https://docs.aws.amazon.com/glue/latest/dg/monitor-tags.html)
	// in the developer guide.
	Tags map[string]string

	// The timeout of the task run for this transform in minutes. This is the maximum
	// time that a task run for this transform can consume resources before it is
	// terminated and enters TIMEOUT status. The default is 2,880 minutes (48 hours).
	Timeout *int32

	// The encryption-at-rest settings of the transform that apply to accessing user
	// data. Machine learning transforms can access user data encrypted in Amazon S3
	// using KMS.
	TransformEncryption *types.TransformEncryption

	// The type of predefined worker that is allocated when this task runs. Accepts a
	// value of Standard, G.1X, or G.2X.
	//   - For the Standard worker type, each worker provides 4 vCPU, 16 GB of memory
	//   and a 50GB disk, and 2 executors per worker.
	//   - For the G.1X worker type, each worker provides 4 vCPU, 16 GB of memory and a
	//   64GB disk, and 1 executor per worker.
	//   - For the G.2X worker type, each worker provides 8 vCPU, 32 GB of memory and a
	//   128GB disk, and 1 executor per worker.
	// MaxCapacity is a mutually exclusive option with NumberOfWorkers and WorkerType .
	//   - If either NumberOfWorkers or WorkerType is set, then MaxCapacity cannot be
	//   set.
	//   - If MaxCapacity is set then neither NumberOfWorkers or WorkerType can be set.
	//   - If WorkerType is set, then NumberOfWorkers is required (and vice versa).
	//   - MaxCapacity and NumberOfWorkers must both be at least 1.
	WorkerType types.WorkerType

	noSmithyDocumentSerde
}

type CreateMLTransformOutput struct {

	// A unique identifier that is generated for the transform.
	TransformId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCreateMLTransformMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpCreateMLTransform{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpCreateMLTransform{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addCreateMLTransformResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpCreateMLTransformValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCreateMLTransform(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addendpointDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opCreateMLTransform(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "glue",
		OperationName: "CreateMLTransform",
	}
}

type opCreateMLTransformResolveEndpointMiddleware struct {
	EndpointResolver EndpointResolverV2
	BuiltInResolver  builtInParameterResolver
}

func (*opCreateMLTransformResolveEndpointMiddleware) ID() string {
	return "ResolveEndpointV2"
}

func (m *opCreateMLTransformResolveEndpointMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	if awsmiddleware.GetRequiresLegacyEndpoints(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	if m.EndpointResolver == nil {
		return out, metadata, fmt.Errorf("expected endpoint resolver to not be nil")
	}

	params := EndpointParameters{}

	m.BuiltInResolver.ResolveBuiltIns(&params)

	var resolvedEndpoint smithyendpoints.Endpoint
	resolvedEndpoint, err = m.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to resolve service endpoint, %w", err)
	}

	req.URL = &resolvedEndpoint.URI

	for k := range resolvedEndpoint.Headers {
		req.Header.Set(
			k,
			resolvedEndpoint.Headers.Get(k),
		)
	}

	authSchemes, err := internalauth.GetAuthenticationSchemes(&resolvedEndpoint.Properties)
	if err != nil {
		var nfe *internalauth.NoAuthenticationSchemesFoundError
		if errors.As(err, &nfe) {
			// if no auth scheme is found, default to sigv4
			signingName := "glue"
			signingRegion := m.BuiltInResolver.(*builtInResolver).Region
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)

		}
		var ue *internalauth.UnSupportedAuthenticationSchemeSpecifiedError
		if errors.As(err, &ue) {
			return out, metadata, fmt.Errorf(
				"This operation requests signer version(s) %v but the client only supports %v",
				ue.UnsupportedSchemes,
				internalauth.SupportedSchemes,
			)
		}
	}

	for _, authScheme := range authSchemes {
		switch authScheme.(type) {
		case *internalauth.AuthenticationSchemeV4:
			v4Scheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4)
			var signingName, signingRegion string
			if v4Scheme.SigningName == nil {
				signingName = "glue"
			} else {
				signingName = *v4Scheme.SigningName
			}
			if v4Scheme.SigningRegion == nil {
				signingRegion = m.BuiltInResolver.(*builtInResolver).Region
			} else {
				signingRegion = *v4Scheme.SigningRegion
			}
			if v4Scheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4Scheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)
			break
		case *internalauth.AuthenticationSchemeV4A:
			v4aScheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4A)
			if v4aScheme.SigningName == nil {
				v4aScheme.SigningName = aws.String("glue")
			}
			if v4aScheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4aScheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, *v4aScheme.SigningName)
			ctx = awsmiddleware.SetSigningRegion(ctx, v4aScheme.SigningRegionSet[0])
			break
		case *internalauth.AuthenticationSchemeNone:
			break
		}
	}

	return next.HandleSerialize(ctx, in)
}

func addCreateMLTransformResolveEndpointMiddleware(stack *middleware.Stack, options Options) error {
	return stack.Serialize.Insert(&opCreateMLTransformResolveEndpointMiddleware{
		EndpointResolver: options.EndpointResolverV2,
		BuiltInResolver: &builtInResolver{
			Region:       options.Region,
			UseDualStack: options.EndpointOptions.UseDualStackEndpoint,
			UseFIPS:      options.EndpointOptions.UseFIPSEndpoint,
			Endpoint:     options.BaseEndpoint,
		},
	}, "ResolveEndpoint", middleware.After)
}