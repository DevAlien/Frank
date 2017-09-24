import React, {Component} from 'react';
import logo from './logo.svg';
import './App.css';
import config from './config';
import io from 'socket.io-client';
import PropTypes from 'prop-types';
import VoiceRecognition from './services/VoiceRecognition';
import CSSTransition from 'react-transition-group/CSSTransition'
import TransitionGroup from 'react-transition-group/TransitionGroup'
import asyncComponent from './components/AsyncComponent';

import {BrowserRouter as Router, Route, Link, Redirect, Switch} from 'react-router-dom'

import Header from './components/Header';

const socket = io(config.frankServerScoket, {transports: ['websocket']});
class App extends Component {
  static childContextTypes = {
    socket: PropTypes.object
  };
  static contextTypes = {
    router: PropTypes.object.isRequired
  };

  constructor(props, context) {
    super(props, context);
  }

  getChildContext() {
    return {socket: socket};
  }

  render() {
    console.log(this.context.router)
    const match = this.context.router.route.match;
    const location = this.context.router.route.location;
    const key = location.pathname.split('/')[1] || '/'

    const timeout = 300
    return (
      <div>
        <Header/>
        <Route exact path="/" render={() => (<Redirect to="/home"/>)}/>

        <TransitionGroup className="page-main">
          <CSSTransition key={key} classNames="fade" timeout={timeout} appear>
            <section className="page-main-inner">
              <Switch location={location}>
                <Route
                  location={location}
                  key={location.key}
                  path="/:component"
                  component={asyncComponent(() => import(`./pages/${key}`))}/>
              </Switch>
            </section>
          </CSSTransition>
        </TransitionGroup>
      </div>
    );
  }
}

// const HSL = ({match: {
//     params
//   }}) => (
//     ())

export default App;
