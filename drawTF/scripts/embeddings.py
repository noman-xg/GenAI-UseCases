from langchain.document_loaders import DirectoryLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.embeddings import OpenAIEmbeddings
from langchain.prompts import PromptTemplate
from langchain.vectorstores import Chroma
from langchain.chains import RetrievalQA
#from langchain.llms import OpenAI
from langchain.chat_models import ChatOpenAI
from prompts import prompt
import argparse
import pandas as pd
import logging
import sys
import openai
import requests
import os
__import__('pysqlite3')
sys.modules['sqlite3'] = sys.modules.pop('pysqlite3')

def main(query,path):
    
    # folder_path = os.path.join ('..', 'vectorStore')
    # vector_store_exists = os.path.isdir (folder_path)
    # print ("Vector Store Status", vector_store_exists)
    # texts = []
    # if not vector_store_exists:
    #     print("VectorStore Does not Exist already, using loadAndSplitDocs")
    #     texts=loadAndSplitDocs(path)
    texts=loadAndSplitDocs(path)
    db, retriever = checkAndQueryVectorStore(texts,query) 
    print("Retriever: ", retriever)
    result = runUserQuery(retriever,query)  
    print("Result: ", result)
    
    # Visualizations for upcoming demo.
    # chunks = [str(text).split('=', 1)[-1][:-1] for text in texts]
    # embDict = db.get(include=['embeddings'])
    # userQueryEmb = visualizeUserQueryEmbedding(query)
    
    # df = pd.DataFrame({'chunks': chunks, 'embeddings': embDict['embeddings']})
    # df.to_csv('output.csv', index=False)
#    print(result)
    return result


def loadAndSplitDocs(path):
  
    loader = DirectoryLoader(path, glob="**/*.tf", use_multithreading=True)
    docs = loader.load()

    text_splitter = RecursiveCharacterTextSplitter(
    # Set a really small chunk size, just to show.
    separators= ['"module ", "resource ", "variable "'],
    chunk_size = 150,
    chunk_overlap  = 20,
    length_function = len,
    add_start_index = True,
    )
    texts = text_splitter.split_documents(docs)
    return texts


def checkAndQueryVectorStore(texts,query):
        
    db = Chroma(persist_directory="./vectorStore", embedding_function=OpenAIEmbeddings())
    embeddingsList = db.get(include=['embeddings'])['embeddings']
    
    if  embeddingsList is None:
        logging.log(level=1, msg="INFO: Creating New VectorStore in the current directory")        
        db = Chroma.from_documents(texts, OpenAIEmbeddings(), persist_directory="vectorStore")        


    retriever = db.as_retriever()
    retriever.search_kwargs["distance_metric"] = "cos"
    retriever.search_kwargs["fetch_k"] = 2
    retriever.search_kwargs["maximal_marginal_relevance"] = False
    retriever.search_kwargs["k"] = 2
    return db,retriever


def runUserQuery(retriever, query):
    
    llm = ChatOpenAI(model="gpt-3.5-turbo-16k-0613",
              temperature=0.5,
              streaming= True, 
              max_tokens = 1000)

    
   # extraction_prompt = PromptTemplate(input_variables=['context', 'question' ],template=prompt)
    extraction_prompt = PromptTemplate(input_variables=['context', 'question'],template="You are an expert at generating Terraform configurations for multiple cloud providers such as AWS,GCP and Azure.Use the following context output either the terraform configuration or list of resources according to the Question. Don't make up any answer, if you don't know the answer just say i dont't know. \n\n{context}\n\nQuestion: {question}\n Helpful Answer:")
    kwargs = {"prompt": extraction_prompt}
    print(extraction_prompt)
    qa = RetrievalQA.from_chain_type(llm=llm, chain_type="stuff", chain_type_kwargs=kwargs,retriever=retriever) 

  #  print("Running Query: ", qa.run(query))
   # print("GPT-Generated Query: ", query)
    return qa.run(query)



# def visualizeUserQueryEmbedding(query):
#     db = Chroma.from_texts([query], OpenAIEmbeddings())
#     #print(db.get(include=['embeddings']))
#     #print("\n\n", len(db.get(include=['embeddings'])['embeddings'][0]))
#     return db.get(include=['embeddings'])['embeddings']



if __name__ == "__main__":
    # Step 2: Create an ArgumentParser object
    parser = argparse.ArgumentParser(description="Script to process a query")

    # Step 3: Define the argument 'query'
    parser.add_argument("query", type=str, help="The query to be processed")
    parser.add_argument("path", type=str, help="Path to the directory containing all documents.")
    # Step 4: Parse the arguments
    args = parser.parse_args()

    #Call the main function with the 'query' argument
    main(args.query,args.path)
