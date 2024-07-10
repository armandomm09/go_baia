import json
import os
import firebase_admin
from firebase_admin import credentials, firestore

# Inicializa la app de Firebase
cred = credentials.Certificate('/Users/armandoo/Downloads/baia-1df5a-firebase-adminsdk-19h0s-4b79382c5c.json')
firebase_admin.initialize_app(cred)
db = firestore.client()

# Define la colección y los documentos
collection_name = 'El Sabor de Tlaxcala'
subcollection_name = 'conversations'
user_id = '5212223201384@c.us'  # ID del usuario (reemplázalo por el ID real)
messages_subcollection = 'rawConversation'
context_prompt_doc = 'Context promt'

# Archivos de salida
output_file = '/Users/armandoo/PROGRA/Golang/go_baia/training/fine_tuning_data.json'
output_file_jsonl = '/Users/armandoo/PROGRA/Golang/go_baia/training/fineTuningBaia.jsonl'

# Leer el contexto del prompt
context_ref = db.collection(collection_name).document(subcollection_name).collection(user_id).document('messages').collection(messages_subcollection).document(context_prompt_doc)
context_doc = context_ref.get()

if context_doc.exists:
    context_data = context_doc.to_dict()
    print("Context prompt encontrado")
else:
    print('No se encontró el documento de contexto del prompt.')

# Función para agregar datos al archivo JSON
def append_to_json_file(data, filename):
    if os.path.isfile(filename):
        with open(filename, 'a', encoding='utf-8') as file:
            file.write('\n' + json.dumps(data, ensure_ascii=False))
    else:
        with open(filename, 'w', encoding='utf-8') as file:
            file.write(json.dumps(data, indent=4, ensure_ascii=False))
    print("Documento agregado")

# Función para convertir JSON a JSONL
def convert_json_to_jsonl(json_filename, jsonl_filename):
    with open(json_filename, 'r', encoding='utf-8') as json_file, open(jsonl_filename, 'w', encoding='utf-8') as jsonl_file:
        for line in json_file:
            jsonl_file.write(line.strip() + '\n')
    print(f"Archivo convertido a {jsonl_filename}")

# Procesar todas las conversaciones
conversations_ref = db.collection(collection_name).document(subcollection_name).collection(user_id).document('messages').collection(messages_subcollection)
counter = 1
messages = []

if context_doc.exists:
    messages.append({'role': context_data['role'], 'content': context_data['content']})

for message_doc in conversations_ref.stream():
    if message_doc.id == context_prompt_doc:
        continue
    print("Counter ", counter)

    message_data = message_doc.to_dict()
    messages.append({
        'role': message_data['role'],
        'content': message_data['content']
    })
    counter += 1

conversation_data = {'messages': messages}
append_to_json_file(conversation_data, output_file)

# Convertir el archivo JSON a JSONL
convert_json_to_jsonl(output_file, output_file_jsonl)
