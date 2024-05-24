# video-encoder: A Golang Video Encoding Tool

## Setting up environment

To execute in development environment, do the following:

- Duplicate `.env.example` to `.env`
- Execute docker-compose up -d
- Open Rabbitmq admin and create an exchange of type`fannout`. It will be a `Dead Letter Exchange` to receive not processed messages.
- Create a `Dead Letter Queue` and bind it to `Dead Letter Exchange`. No need for a routing_key.
- In `.env` pass the `Dead Letter Exchange` name in: `RABBITMQ_DLX`
- Create an account in GCP that has the permission to store in GCS. Download the json file with the credentials and save it in the root of the project with the name: `bucket-credential.json`

## Executing

Run `make server` inside the container:

```
docker exec encoder-new2_app_1 make server
```

`microsservico-enconder_app_1` is the container name generated by docker-compose.

## Message format to send to the encoder

```
{
  "resource_id": "my-resource-id-can-be-a-uuid-type",
  "file_path": "convite.mp4"
}
```

- `resource_id`: Video ID, a string.
- `file_path`: Full path for the mp4 inside the bucket.

## Message format returned by the encoder

### Success response

For each video processed, the encoder will send to an exchange (setup in .env) the following JSON:

```
{
    "id":"bbbdd123-ad05-4dc8-a74c-d63a0a2423d5",
    "output_bucket_path":"codeeducationtest",
    "status":"COMPLETED",
    "video":{
        "encoded_video_folder":"b3f2d41e-2c0a-4830-bd65-68227e97764f",
        "resource_id":"aadc5ff9-0b0d-13ab-4a40-a11b2eaa148c",
        "file_path":"convite.mp4"
    },
    "Error":"",
    "created_at":"2020-05-27T19:43:34.850479-04:00",
    "updated_at":"2020-05-27T19:43:38.081754-04:00"
}
```

`encoded_video_folder` is the folder of the enconded video.

### Error response

For each process that failed, a JSON will be returned:

```
{
    "message": {
        "resource_id": "aadc5ff9-010d-a3ab-4a40-a11b2eaa148c",
        "file_path": "convite.mp4"
    },
    "error":"Motivo do erro"
}
```

Also, the encoder will send to a dead letter exchange the original message that failed to process.
Setup the desired DLX in the .env file parameter: `RABBITMQ_DLX`
