````markdown id="r6n1dw"
# Go CLI Download Manager

A small concurrency project built to practice:

- goroutines
- channels
- semaphores
- concurrent file downloads
- limiting concurrency

This program downloads several images concurrently and stores them in a `downloads/` folder.

---

## Concepts Practiced

- Goroutines
- Channels
- WaitGroups
- Buffered channels
- Semaphore pattern
- Concurrent HTTP requests
- File creation and saving

---

## How It Works

1. A list of image URLs is created.
2. A `downloads/` folder is created if it does not already exist.
3. A goroutine is started for each URL.
4. A semaphore limits the number of simultaneous downloads to 3.
5. Each download:
   - fetches the image
   - saves it to the `downloads/` folder
   - reports success or failure through a channel
6. The main goroutine prints the results as downloads finish.

---

## Project Structure

```text
.
├── downloads/
├── main.go
└── README.md
```
````

---

## Example URLs

The project uses placeholder images from Picsum Photos:

```go id="o3e3ur"
urls := []string{
	"https://picsum.photos/200/300",
	"https://picsum.photos/300/300",
	"https://picsum.photos/400/300",
	"https://picsum.photos/500/300",
	"https://picsum.photos/600/300",
}
```

---

## Why a Semaphore Is Needed

Without a limit, every download would start at the same time.

This project limits the number of active downloads to 3 using:

```go id="l1qf4q"
semaphore := make(chan struct{}, 3)
```

The buffered channel acts like a gate:

- If fewer than 3 downloads are running, a new one may start.
- If 3 downloads are already running, the next one waits.

Example flow:

```text
Download 1 starts
Download 2 starts
Download 3 starts

One finishes
→ Download 4 starts

Another finishes
→ Download 5 starts
```

---

## How the Semaphore Works

Before a download begins:

```go id="lfp0h4"
semaphore <- struct{}{}
```

This takes one slot.

When the download finishes:

```go id="t6lmyq"
<-semaphore
```

This frees the slot for another download.

---

## Unique Filenames

The URLs from Picsum Photos do not contain useful file names.

For example:

```text
https://picsum.photos/200/300
```

would normally become:

```text
300
```

using `filepath.Base(url)`.

That would cause all files to overwrite each other.

Instead, the project creates unique names:

```go id="e0rymn"
filename := fmt.Sprintf("image_%d.jpg", id)
```

Resulting files:

```text
downloads/
├── image_1.jpg
├── image_2.jpg
├── image_3.jpg
├── image_4.jpg
└── image_5.jpg
```

---

## Running the Program

Make sure Go is installed.

Run:

```bash id="cm9js3"
go run .
```

---

## Expected Output

The order may change each run because downloads happen concurrently.

Example:

```text id="g1ncf0"
Starting download: image_1.jpg
Starting download: image_2.jpg
Starting download: image_3.jpg

Finished downloading: downloads\image_2.jpg
Finished downloading: downloads\image_1.jpg

Starting download: image_4.jpg
Starting download: image_5.jpg

Finished downloading: downloads\image_3.jpg
Finished downloading: downloads\image_4.jpg
Finished downloading: downloads\image_5.jpg

All downloads completed.
```

Notice that:

- Only 3 downloads are active at once.
- A new download begins only after another one finishes.
- Downloads do not necessarily finish in the same order they started.

---

## Files Saved

After running the project, the `downloads/` folder should contain:

```text id="h1lylt"
downloads/
├── image_1.jpg
├── image_2.jpg
├── image_3.jpg
├── image_4.jpg
└── image_5.jpg
```

---

## What You Learned

By completing this project, you practiced:

- launching many goroutines
- limiting concurrency
- synchronizing goroutines with a `WaitGroup`
- communicating results through channels
- downloading and saving files concurrently

---

## Possible Improvements

You can extend the project by:

- showing download percentages
- saving the original content type and extension
- adding timeouts:

```go id="h4r9of"
client := http.Client{
	Timeout: 10 * time.Second,
}
```

- retrying failed downloads
- allowing URLs to be passed through command-line arguments
- adding colored terminal output
- showing how many downloads are currently active
- creating a progress bar

```

```
