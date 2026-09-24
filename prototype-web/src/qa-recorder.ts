export class CanvasRecorder {
  private recorder?: MediaRecorder;
  private stream?: MediaStream;
  private timeout?: number;
  private chunks: Blob[] = [];
  private recording?: Blob;
  private started = false;
  private failure = "";

  get supported(): boolean {
    return "MediaRecorder" in window && "captureStream" in HTMLCanvasElement.prototype;
  }

  get ready(): boolean {
    return this.recording !== undefined;
  }

  get state(): string {
    if (this.failure) return `error:${this.failure}`;
    if (this.ready) return "ready";
    if (this.recorder) return this.recorder.state;
    return this.started ? "stopped" : "idle";
  }

  start(canvas: HTMLCanvasElement, framesPerSecond = 60): boolean {
    if (!this.supported || this.recorder?.state === "recording") return false;
    try {
      const stream = canvas.captureStream(framesPerSecond);
      this.stream = stream;
      const preferred = "video/webm;codecs=vp9";
      const mimeType = MediaRecorder.isTypeSupported(preferred) ? preferred : "video/webm";
      this.chunks = [];
      this.recording = undefined;
      this.failure = "";
      this.recorder = new MediaRecorder(stream, { mimeType, videoBitsPerSecond: 12_000_000 });
      this.recorder.addEventListener("dataavailable", (event) => {
        if (event.data.size > 0) this.chunks.push(event.data);
      });
      this.recorder.addEventListener("stop", () => {
        this.recording = new Blob(this.chunks, { type: mimeType });
        this.chunks = [];
        this.releaseTracks();
      });
      this.recorder.addEventListener("error", (event) => {
        this.failure = event.error.message;
        this.releaseTracks();
      });
      this.recorder.start(250);
      this.timeout = window.setTimeout(() => { void this.stop(); }, 180_000);
      this.started = true;
      return true;
    } catch (error) {
      this.releaseTracks();
      this.failure = error instanceof Error ? error.message : String(error);
      return false;
    }
  }

  async stop(): Promise<Blob | undefined> {
    if (!this.recorder || this.recorder.state === "inactive") return this.recording;
    const recorder = this.recorder;
    await new Promise<void>((resolve) => {
      recorder.addEventListener("stop", () => resolve(), { once: true });
      recorder.stop();
    });
    return this.recording;
  }

  private releaseTracks(): void {
    if (this.timeout !== undefined) window.clearTimeout(this.timeout);
    this.stream?.getTracks().forEach((track) => track.stop());
    this.stream = undefined;
  }

  dispose(): void {
    if (this.recorder?.state === "recording") this.recorder.stop();
    this.releaseTracks();
  }

  async asDataUrl(): Promise<string | undefined> {
    if (!this.recording) return undefined;
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.addEventListener("load", () => resolve(String(reader.result)), { once: true });
      reader.addEventListener("error", () => reject(reader.error), { once: true });
      reader.readAsDataURL(this.recording!);
    });
  }

}
