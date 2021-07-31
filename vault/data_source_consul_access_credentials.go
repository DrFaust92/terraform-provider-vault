package vault

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/vault/api"
)

func consulAccessCredentialsDataSource() *schema.Resource {
	return &schema.Resource{
		Read: consulCredentialsDataSourceRead,

		Schema: map[string]*schema.Schema{
			"backend": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Consul Secret Backend to read credentials from.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Consul Secret Role to read credentials from.",
			},

			"lease_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lease identifier assigned by vault.",
			},

			"lease_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Lease duration in seconds relative to the time in lease_start_time.",
			},

			"lease_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which the lease was read, using the clock of the system where Terraform was running",
			},

			"lease_renewable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True if the duration of this lease can be extended through renewal.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lease identifier assigned by vault.",
			},
			// "ttl": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "User specified Time-To-Live for the STS token. Uses the Role defined default_sts_ttl when not specified",
			// },
		},
	}
}

func consulCredentialsDataSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	backend := d.Get("backend").(string)
	credType := d.Get("type").(string)
	role := d.Get("role").(string)
	path := backend + "/" + credType + "/" + role

	// arn := d.Get("role_arn").(string)
	// // If the ARN is empty and only one is specified in the role definition, this should work without issue
	// data := map[string][]string{
	// 	"role_arn": {arn},
	// }

	// if v, ok := d.GetOk("ttl"); ok {
	// 	data["ttl"] = []string{v.(string)}
	// }

	log.Printf("[DEBUG] Reading %q from Vault", path)
	secret, err := client.Logical().Read(path)
	if err != nil {
		return fmt.Errorf("error reading from Vault: %s", err)
	}
	log.Printf("[DEBUG] Read %q from Vault", path)

	if secret == nil {
		return fmt.Errorf("no role found at path %q", path)
	}

	// accessKey := secret.Data["access_key"].(string)
	// secretKey := secret.Data["secret_key"].(string)
	// var securityToken string
	// if secret.Data["security_token"] != nil {
	// 	securityToken = secret.Data["security_token"].(string)
	// }

	d.SetId(secret.LeaseID)
	d.Set("token", secret.Data["token"])
	// d.Set("secret_key", secret.Data["secret_key"])
	// d.Set("security_token", secret.Data["security_token"])
	d.Set("lease_id", secret.LeaseID)
	d.Set("lease_duration", secret.LeaseDuration)
	d.Set("lease_start_time", time.Now().Format(time.RFC3339))
	d.Set("lease_renewable", secret.Renewable)

	// awsConfig := &aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(accessKey, secretKey, securityToken),
	// 	HTTPClient:  cleanhttp.DefaultClient(),
	// }

	// region := d.Get("region").(string)
	// if region != "" {
	// 	awsConfig.Region = &region
	// }

	// sess, err := session.NewSession(awsConfig)
	// if err != nil {
	// 	return fmt.Errorf("error creating AWS session: %s", err)
	// }

	// iamconn := iam.New(sess)
	// stsconn := sts.New(sess)

	// // Different types of AWS credentials have different behavior around consistency.
	// // See https://www.vaultproject.io/docs/secrets/aws/index.html#usage for more.
	// if credType == "sts" {
	// 	// STS credentials are immediately consistent. Let's ensure they're working.
	// 	log.Printf("[DEBUG] Checking if AWS sts token %q is valid", secret.LeaseID)
	// 	if _, err := stsconn.GetCallerIdentity(&sts.GetCallerIdentityInput{}); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	// // Other types of credentials are eventually consistent. Let's check credential
	// // validity and slow down to give credentials time to propagate before we return
	// // them. We'll wait for at least 5 sequential successes before giving creds back
	// // to the user.
	// sequentialSuccesses := 0

	// // validateCreds is a retry function, which will be retried until it succeeds.
	// validateCreds := func() *resource.RetryError {
	// 	log.Printf("[DEBUG] Checking if AWS creds %q are valid", secret.LeaseID)
	// 	if _, err := iamconn.GetUser(nil); err != nil && isAWSAuthError(err) {
	// 		sequentialSuccesses = 0
	// 		log.Printf("[DEBUG] AWS auth error checking if creds %q are valid, is retryable", secret.LeaseID)
	// 		return resource.RetryableError(err)
	// 	} else if err != nil {
	// 		log.Printf("[DEBUG] Error checking if creds %q are valid: %s", secret.LeaseID, err)
	// 		return resource.NonRetryableError(err)
	// 	}
	// 	sequentialSuccesses++
	// 	log.Printf("[DEBUG] Checked if AWS creds %q are valid", secret.LeaseID)
	// 	return nil
	// }

	// start := time.Now()
	// for sequentialSuccesses < sequentialSuccessesRequired {
	// 	if time.Since(start) > sequentialSuccessTimeLimit {
	// 		return fmt.Errorf("unable to get %d sequential successes within %.f seconds", sequentialSuccessesRequired, sequentialSuccessTimeLimit.Seconds())
	// 	}
	// 	if err := resource.Retry(retryTimeOut, validateCreds); err != nil {
	// 		return fmt.Errorf("error checking if credentials are valid: %s", err)
	// 	}
	// }

	// log.Printf("[DEBUG] Waiting an additional %.f seconds for new credentials to propagate...", propagationBuffer.Seconds())
	// time.Sleep(propagationBuffer)
	return nil
}
