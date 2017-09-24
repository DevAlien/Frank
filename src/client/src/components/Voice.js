import React, { Component } from 'react';
import PropTypes from 'prop-types';
import VoiceRecognition from '../services/VoiceRecognition';
import FrankSleep from '../icons/FrankSleep';
import FrankListen from '../icons/FrankListen';
import FrankProgress from '../icons/FrankProgress';

const sleeping = 'sleeping';
const listening = 'listening';
const processing = 'processing';

const stopListeningTimeout = 1000 * 30;
const wakeUps = ['frank', 'svegliati'];
const sleeps = ['stai zitto', 'stai zitta', 'spegniti', 'vai a dormire frank'];
export default class Voice extends Component {

  static contextTypes = {
    socket: PropTypes.object.isRequired
  };

  constructor(props, context) {
    super(props, context);
    
    const options = {
      onText: this.onText
    };
    this.voiceRecognition = new VoiceRecognition(options);
    this.state = {
      status: sleeping
    }
  }

  componentDidMount() {
    this.voiceRecognition.start();
  }

  onText = (text) => {
    text = text.toLowerCase();


    if (this.state.status !== sleeping) {
      if (sleeps.indexOf(text) !== -1) {
        this.setState({status: sleeping});
        this.clearSleepingTimeout();
      } else {
        this.setState({status: processing});
        return this.context.socket.emit("text", text, (result, value) => {
          this.setState({status: listening});
          if (result === true) {
            this.startSleepingTimeout();
          }
          console.log('r', result)
          console.log('v', value)
        });
      }
    }

    if (wakeUps.indexOf(text) !== -1) {
      this.setState({status: listening});
      this.startSleepingTimeout();
    }

    
  }

  startSleepingTimeout() {
    if (this.listeningTimeout) {
      this.clearSleepingTimeout()
    }

    this.listeningTimeout = setTimeout(() => {
      this.setState({status: sleeping});
    }, stopListeningTimeout);
  }

  clearSleepingTimeout() {
    clearTimeout(this.listeningTimeout);
    this.listeningTimeout = false;
  }
  st = () => {
    console.log('start')
    //this.voiceRecognition.abort()
    this.voiceRecognition.start()
  }
  render() {
    let icon = <FrankSleep onClick={() => this.st()}className="nav-icon"/>;
    switch(this.state.status) {
      case listening:
        icon = <FrankListen className="nav-icon icon-pulse"/>;
        break;
      case processing:
        icon = <FrankProgress  className="nav-icon icon-spin"/>;
        break;
    }

    return <div className="frank-status">{icon}</div>;
  }
}