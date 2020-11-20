import './header.css'

import {Link} from 'react-router-dom';
import React from 'react';
import { connect } from "react-redux";
import { logout } from "../../actions/auth";

class Header extends React.Component{
    constructor(props) {
        super(props);
        this.logout = this.logout.bind(this);
    }
    logout() {
        this.props.dispatch(logout());
    }

    render(){
        const {isLoggedIn} = this.props;
        return(
            <div className="bg-purple-900 w-full px-4 py-2 flex justify-center">
                <div className="max-w-3xl w-full ">
                    <div className="flex items-center justify-between text-purple-200">
                        <Link to='/'>
                            <div className="text-purple-100  text-2xl font-mono hover:text-purple-300 cursor-pointer select-none">EasySubmit</div>
                        </Link>
                        <div className="flex">     
                            {
                                isLoggedIn ? (
                                    <div>
                                        <Link to="/profile">
                                            <button type="button" className="focus:outline-none active:border-purple-900 hover:text-purple-300 hover:border-purple-300 w-16 h-8 border-2 border-purple-100 rounded-lg  text-xs mr-2 font-thin p-1 ">My profile</button>
                                        </Link>
                                        <Link to='/'>
                                            <button type="button" onClick={this.logout} className="focus:outline-none w-16 h-8 border-2 hover:text-purple-300 hover:border-purple-300 border-purple-100 rounded-lg  text-xs font-thin p-1">Sign out</button>                    
                                        </Link>
                                    </div>
                                ) : 
                                (
                                    <div>
                                        <Link to="/signin">
                                            <button type="button" className="focus:outline-none active:border-purple-900 hover:text-purple-300 hover:border-purple-300 w-16 h-8 border-2 border-purple-100 rounded-lg  text-xs mr-2 font-thin p-1 ">Sign in</button>
                                        </Link>
                                        <Link to='/signup'>
                                            <button type="button" className="focus:outline-none w-16 h-8 border-2 hover:text-purple-300 hover:border-purple-300 border-purple-100 rounded-lg  text-xs font-thin p-1">Sign up</button>                    
                                        </Link>
                                    </div> 
                                )   
                            }                   
                                
                        </div>
                    </div>
                </div>
            </div>
        )
    }
    
}
function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    return {
      isLoggedIn
    };
  }
  
export default connect(mapStateToProps)(Header);

