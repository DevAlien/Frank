class VoiceRecognition {
  constructor(options) {
    this.options = options || {};
    const BrowserSpeechRecognition =
    window.SpeechRecognition ||
    window.webkitSpeechRecognition ||
    window.mozSpeechRecognition ||
    window.msSpeechRecognition ||
    window.oSpeechRecognition

    this.recognition = new BrowserSpeechRecognition();
    this.recognition.continuous = true;
    this.recognition.interimResults = true;
    this.recognition.onstart = this.onStart;
    this.recognition.onend = this.onEnd;
    this.recognition.onerror = this.onError;
    this.recognition.onresult = this.onResult;
    this.recognition.lang = "it_IT";

    this.text = '';
  }

  onStart = () => {
    if (this.options.onStart) this.options.onStart();
  }

  onEnd = () => {
    if (this.options.onEnd) this.options.onEnd();

    setTimeout(() => {
      this.recognition.start();
    }, 0);
  }

  onError = () => {
    if (this.options.onError) this.options.onError();
  }

  onResult = (event) => {
    if (this.options.onEnd) this.options.onResult(event);

    if (this.timeout) {
      clearTimeout(this.timeout); 
    }

    let interimTranscript = '';
    for (var i = event.resultIndex; i < event.results.length; ++i) {
      if (event.results[i].isFinal) {
        this.text += event.results[i][0].transcript;
      } else {
        interimTranscript += event.results[i][0].transcript;
      }
    }

    if (this.text.length > 0) {
      console.log('return final')
      this.returnText(this.text);
    } else {
      console.warn('intermediate', interimTranscript)
      this.timeout = setTimeout(() => {console.log('return inter'); this.returnText(interimTranscript)}, 300);
    }
  }

  returnText(text) {
    console.log('send', text)
    if (this.options.onText) this.options.onText(text);

    this.text = '';
    this.recognition.abort();
    
  }

  start() {
    this.recognition.start();
  }

  stop() {
    this.recognition.abort();
  }
}

export default VoiceRecognition;