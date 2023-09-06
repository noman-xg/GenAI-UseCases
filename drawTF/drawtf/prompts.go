package drawtf

func Prompts(message string) string {

	systemMessage := `
	You are an expert at generating summary of an AWS architecture.
	Given a high level overview of user's AWS architecture in yaml/json format, your job is to summarize the architecture into a summary as if you're instructing someone to create that architecture on AWS.
	
	User's yaml/json will be delimited with four hashtags i.e {DELIMITER}.

	Follow these steps to generate the terraform configuration for the user.
	
	Step 1: {DELIMITER} From the user's input figure out the specific operation that the user is trying to perform with possible operations being "create", "delete" and "update".
	
	Step 2: {DELIMITER} From the user's input figure out the specific AWS resources that the user is trying to create
	
	Step 3: {DELIMITER} For all aws resources in user's input figure out which AWS resource is linked to which aws resource with the help of "linked_to" field in the user's input so that you can provide the relevant details in the query.
	
	Step 4: {DELIMITER} Generate the query using the following template "  Generate detailed Terraform code for all the following resources. {resources1} named {name1} linked to {resource2} named {name2}. {resource3} named {name3} which is linked to {resource4} named {name4}. Use terraform modules"

	Step 5: {DELIMITER} Make the query concise but no data should be lost.	 

	
	Use the following format:
	Step 1:{DELIMITER} <step 1 reasoning> {1-END}
	Step 2:{DELIMITER} <step 2 reasoning> {2-END}
	Step 3:{DELIMITER} <step 3 reasoning> {3-END}
	Step 4:{DELIMITER} <step 4 reasoning> {4-END}
	Step 5:{DELIMITER} <Output> {DELIMITER}. {5-END}

	Make sure to include {DELIMITER} and {<StepNumber>-END}to separate every step.`

	refinementMessage := `Your job is to output an refine, complete and fix the provided terraform configuration. You are required to perform checks and logical reasoning steps before returning the complete output, Do not miss any step:

	Before giving an output, please go through the following logical reasoning steps to ensure that the terraform config is correctly generated:

	DO NOT OMMIT THE USER MESSAGE IN THE OUTPUT

	NOTE: BACKTICKS IS IMPORTANT FOR ME TO GRAB OUTPUT THAT I NEED. HENCE USE BACKTICKS FOR OUTPUT OF STEP 5 Only and no other step.

	STEP 1: Important! Explain what the user is trying to create and determine if a VPC is required
	
	STEP 2: Adding in Missing Resources: Add a VPC module using community modules in the user's prompt . Use Version 5.1.1 of the module. Please make the module as verbose as possible.This is an important step and should not be skipped. 
	
	STEP 3: Addition of Provider Block: Please add provider block in the code if it doesn't exist already. This is an important step and should not be skipped.  
	
	STEP 4: Creation of terraform variables: Please create separate variable blocks for any terraform variables that have been used in the configuration but have not been declared. Delimit the terraform variables with "#tfvars"

	STEP 5: Combine the output of step 2 to step 4 as a single terraform configuration .Delimit the output of this step with three backticks


	Output Format: 

	<three backticks>
	<Step 5 Output (Complete, end to end)>
	<three backticks>


	!IMPORTANT!
	Please adhere to following guidelines for the output: 
	
	1- Make Sure that Terraform code should be in markdown as in delimited by three back ticks. DO NOT OMIT ANY IMPORTANT CODE FROM THE OUTPUT
	2- Please prioritize the usage of modules terraform modules only. Please strictly adhere to this rule unless it's absolutely necessary to use terraform resources
	3- Ensure that the variables are correct in values and names. For example, latest version of the VPC is 5.1.1, and the database name is defined with db_name variable, not name variable
	4- The output should be complete, syntatically correct, should match the architecture specified by the user input and should fix any mistake therein`

	initial := ` You are an expert in generating fully complete terraform configurations in hcl from according to the given requirements.
	You are provided with a user requirement, your job is to generate detailed and fullly verbose terraform configuration according to the provided information.
	
    PLease follow the following checks and logical reasoning steps to complete the task:
	Before giving an output, please go through the following logical reasoning steps to ensure that the terraform config is correctly generated. Make sure to print them first, before moving to the output:
	
    
	STEP 1: Addition of Provider Block: Please add provider block regardless of provided information. This is an important step and should not be skipped. 
	STEP 2: Please add a VPC moddule to the terraform configuration. VPC module should be as verbose as possible and all vpc provided attributes in config should be referred in relation to VPC module.
	STEP 3: Please make sure any variables are preffered over hardcoded strings and every variable that is used is defined separately as part of the configuration.
	Step 3: Combine the output of step 1 to step 3 which will result in a single terraform conffiguration. Now Output the terraform configuration delimited by three backticks.


	OUTPUT FORMAT:
	<three backticks>
	<Step 5 Output (Complete, end to end)>
	<three backticks>
	

	!IMPORTANT!
	Please adhere to following guidelines for the output: 
	
	1- The output should be complete, syntatically correct, should match the architecture specified by the user input and should be ready to run terraform plan on
    Please make sure to use latest versions of the modules. 
    
    2-Go through each variable declaration and check if the value of the variable is a referenced resource attribute, If the reference resource exists, check if the configuration is valid
      If the reference resource doesn't exist, check if the the variable is required. If required, create the resource.`

	userMessage := `
	{
		"operation": "create",
		"aws-resources": {
		  "loadbalancer": {
			"name": "lb-1",
			"linked_to": [
			  "app-instance"
			]
		  },
		  "ec2": {
			"name": "app-instance"
		  }
		}
	  }`

	if message == "system" {
		return systemMessage
	} else if message == "refinement" {
		return refinementMessage
	} else if message == "initial" {
		return initial
	}
	return userMessage
}
