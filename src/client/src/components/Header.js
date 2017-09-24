import React, {Component} from 'react'
import {NavLink} from 'react-router-dom'

import Cog from '../icons/Cog';
import Home from '../icons/Home';
import Apps from '../icons/Apps';
import Voice from './Voice';

export default class Header extends Component {
  constructor(props) {
    super(props);

    this.state = {
      showApps: false
    }
    this.r = null;
  }
  
  clickApps = () => {
    this.setState({showApps: !this.state.showApps});
  }

  render() {
    return (
      <div>
        <header className="header-container">
          <div className="header-left">
            <NavLink to="/home" className="nav-link"><Home className="nav-icon"/></NavLink>
            <NavLink to="/settings" className="nav-link"><Cog className="nav-icon"/></NavLink>
          </div>
          <div className="header-center"><Apps className="nav-icon apps-icon" onClick={this.clickApps}/></div>
          <div className="header-right"><Voice ref={(r) => this.voice = r}/></div>

        </header>
        {this.state.showApps &&
        <div className="launchpad shown">
          <div className="content">
            <div>

              <NavLink to="/home" className="icon" onClick={this.clickApps}>
                <Home className="launchpad-icon"/>
                <span>Home</span>
              </NavLink>
              <NavLink to="/settings" className="icon" onClick={this.clickApps}>
                <Cog className="launchpad-icon"/>
                <span>Settings</span>
              </NavLink>
            </div>
          </div>
        </div>
        }
      </div>
    );
  }
}