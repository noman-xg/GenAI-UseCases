prompt1 = """
You are an expert at generating Terraform configurations.
Use the provided context to output detailed terraform configuration. 

Please adhere striclty to these guidelines:

- You can only add or modify resource/module attributes.
- You don't have permission to remove any attribute

Don't make up any answer, if you don't know the answer just say i dont't know. MAKE SURE You only output the terraform code with nothing else. 
Only output terraform without any description.
Terraform code should be delimited by ####
\n\n{context}\n\nQuestion: {question}\n Helpful Answer:"
"""
prompt = """
    You are an generate detailed and fullly verbose terraform configuration from the provided data.
	
    PLease follow the following checks and logical reasoning steps to complete the task:
	Before giving an output, please go through the following logical reasoning steps to ensure that the terraform config is correctly generated. Make sure to print them first, before moving to the output:
	
    
	STEP 1: Addition of Provider Block: Please add provider block regardless of provided information. This is an important step and should not be skipped. 
    STEP 2: When using terraform modules please use a palce-holder for "version" to avoid version conflicts.
	Step 3: Output the terraform configuration delimited by four hashtags.

	!IMPORTANT!
	Please adhere to following guidelines for the output: 
	
	1- The output should be complete, syntatically correct, should match the architecture specified by the user input and should be ready to run terraform plan on
    Please make sure to use latest versions of the modules. 
    
    2-Go through each variable declaration and check if the value of the variable is a referenced resource attribute, If the reference resource exists, check if the configuration is valid
      If the reference resource doesn't exist, check if the the variable is required. If required, create the resource.
    Terraform code should be delimited by ####. CONTEXT : \n{context}\n Question: {question}\nHelpful Answer:
    """