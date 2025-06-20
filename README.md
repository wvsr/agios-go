# 📘 API Documentation

**Backend Stack:** Golang + PostgreSQL  
**API Version:** `v1`  
**Base URL:** `/api/v1/`

---

## ⚙️ Event Types

- `START`
- `END`
- `PLAN`
- `WEB_RESULTS`
- `MARKDOWN_ANSWER`
- `WIDGET`

---

## 📦 Widget Types

- `WEATHER_WIDGET`
- `CRYPTO_WIDGET`
- `NEARBY_PLACES_WIDGET`
- `YT_SUMMARY_WIDGET`
- `SYNTHESIZER_RESULTS`

---

## 📊 Health Check

### `GET /health`

**Description:**  
Used by monitoring tools to check if the server is running.

#### ✅ Response

```http
200 OK
OK
```

---

## 📁 Upload File

### `POST /api/v1/files/upload`

**Description:**  
Upload up to 5 files before submitting a query.

#### 📌 Constraints:

- Max **5 files** per request
- Max **10MB** per file
- Supported MIME types:
  - application/pdf
  - application/x-javascript, text/javascript
  - application/x-python, text/x-python
  - text/plain, text/html, text/css, text/md, text/csv, text/xml, text/rtf
  - image/png, image/jpeg, image/webp, image/heic, image/heif

#### 🔐 Headers

```http
Content-Type: multipart/form-data
```

#### 📤 Request Body

- `files[]`: array of files

#### ✅ Response `200 OK`

```json
[
  {
    "id": "uuid",
    "file_name": "internal.pdf",
    "original_file_name": "user_uploaded.pdf",
    "file_size_bytes": 304920,
    "mime_type": "application/pdf",
    "uploaded_at": "2025-06-20T12:00:00Z",
    "version": "1.0"
  }
]
```

#### ❌ Error Responses

```json
{
  "error": {
    "message": "File size exceeds 10MB limit.",
    "code": "FILE_TOO_LARGE"
  }
}
```

```json
{
  "error": {
    "message": "Maximum 5 files allowed per upload.",
    "code": "MAX_FILE_COUNT_EXCEEDED"
  }
}
```

---

## 🧵 Create Thread (Initial Message)

### `POST /api/v1/threads`

**Description:**  
Create a thread and send the first message.

**Note:** Ensure `GET /api/v1/threads/:threadId` returns no messages before calling this. And if specific file id doesn't exist, it will skip that specific file.

#### 📤 Request Body

```json
{
  "slug": "hello-world-thread",
  "query_text": "What's the weather in Tokyo?",
  "file_ids": ["uuid1", "uuid2"] // optional
}
```

#### ❌ Error Responses

```json
{
  "error": {
    "message": "Maximum 5 files allowed per upload.",
    "code": "MAX_FILE_COUNT_EXCEEDED"
  }
}
```

```json
{
  "error": {
    "message": "File size exceeds 10MB limit.",
    "code": "FILE_TOO_LARGE"
  }
}
```

#### 🔁 Response: **SSE stream**

```
event: START
data: {}

event: PLAN
data: {...}
```

---

## 💬 Add Message to Existing Thread

### `POST /api/v1/threads/:threadId/messages`

**Description:**  
Adds a follow-up message to an existing thread.

#### 📤 Request Body

```json
{
  "query_text": "How about Kyoto?",
  "file_ids": ["uuid3"]
}
```

#### 🔁 Response: **SSE stream**

```
event: START
data: {}

event: MARKDOWN_ANSWER
data: { "chunk": "...", "streaming": true, }

event: WIDGET
data: {
  "version": "1.0",
  "widget_type": "WEATHER_WIDGET",
  "streaming": true,
}
```

#### ❌ Error Example

```json
{
  "error": {
    "message": "Thread not found.",
    "code": "THREAD_NOT_FOUND"
  }
}
```

```json
{
  "error": {
    "message": "Query text exceeds the 1000-word limit.",
    "code": "QUERY_TEXT_TOO_LONG"
  }
}
```

---

## 📜 Get Thread with Messages

### `GET /api/v1/threads/:threadId`

**Description:**  
Returns a thread and its messages.

#### ✅ Response `200 OK`

```json
{
  "id": "uuid",
  "slug": "hello-world-thread",
  "created_at": "2025-06-20T12:00:00Z",
  "updated_at": "2025-06-20T12:15:00Z",
  "version": "1.0",
  "messages": [
    {
      "id": "uuid",
      "query_text": "What's the weather in Tokyo?",
      "response_text": "Currently 33°C and sunny.",
      "event_type": "WIDGET",
      "stream_status": "DONE",
      "meta_data": {
        "widget": {
          "widget_type": "WEATHER_WIDGET",
          "widget_data": {}
        },
        "file_ids": []
      },
      "message_index": 0,
      "created_at": "2025-06-20T12:00:01Z",
      "version": "1.0"
    }
  ]
}
```

#### ❌ Error Example

```json
{
  "error": {
    "message": "Thread not found.",
    "code": "THREAD_NOT_FOUND"
  }
}
```

---

## 🗑️ Delete Thread

### `DELETE /api/v1/threads/:threadId`

**Description:**  
Deletes a thread and all its messages.

#### ✅ Response

```http
204 No Content
```

#### ❌ Error

```json
{
  "error": {
    "message": "Thread not found.",
    "code": "THREAD_NOT_FOUND"
  }
}
```

---

## 🗑️ Delete Message

### `DELETE /api/v1/messages/:messageId`

**Description:**  
Deletes a specific message.

#### ✅ Response

```http
204 No Content
```

#### ❌ Error Example

```json
{
  "error": {
    "message": "Message not found.",
    "code": "MESSAGE_NOT_FOUND"
  }
}
```

---

## 🧩 Event Flow (SSE Stream)

```
event: START
data: {"streaming": true}

event: PLAN
data: { version: string, "cot": "...", "streaming": true }

event: PLAN
data: { version: string, "cot": "...", "streaming": true }

event: WEB_RESULTS
data: { "results": [...], "streaming": true }

event: MARKDOWN_ANSWER
data: { "chunk": "...", "streaming": true, }

event: WIDGET
data: {
  "version": "1.0",
  "widget_type": "SYNTHESIZER_RESULTS",
  "streaming": true,
  widget_data: {}
}

event: END
data: {"streaming": false}
```

---

## 🧠 Widget Response Formats

### 🔹 SYNTHESIZER_RESULTS

```json
{
  "version": "1.0",
  "key_takeaways": [
    {
      "text": "Summary point here",
      "confidence_score": 0.92
    }
  ],
  "related_search_terms": ["machine learning", "ai summary"],
  "short_summary": "This is a short summary.",
  "metrics": [
    {
      "title": "Estimated Time Saved",
      "value": "3 minutes"
    }
  ]
}
```

---

### 🔹 YT_SUMMARY_WIDGET

```json
{
  "version": "1.0",
  "youtube_url": "https://youtube.com/watch?v=abc123"
}
```

---

## 🛠️ Frontend Integration Notes

- Use `EventSource` or `ReadableStream` to handle SSE.
- Always stream in order: `START → PLAN → WIDGET → MARKDOWN_ANSWER`.
- Pre-upload files before sending message.
- On first message: check that thread has no messages.

---

## 🧬 Entity Versioning

Semantic versioning (`"1.0"`, `"1.1"`, etc.) is supported on:

- `UploadFile.version`
- `Thread.version`
- `Message.version`

Use this to track schema evolution across deployments.

---
