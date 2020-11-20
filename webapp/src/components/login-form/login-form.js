import CheckButton from "react-validation/build/button";
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import React from 'react';
import { Redirect } from 'react-router-dom';
import { connect } from "react-redux";
import { login } from "../../actions/auth";

const required = (value) => {
    if (!value) {
      return (
        <div className="border border-rounded bg-red-200 text-red-800 text-center" role="alert">
          This field is required!
        </div>
      );
    }
  };

class LoginForm extends React.Component {
    constructor(props) {
        super(props);
        this.handleLogin = this.handleLogin.bind(this);
        this.onChangeEmail = this.onChangeEmail.bind(this);
        this.onChangePassword = this.onChangePassword.bind(this);
    
        this.state = {
          email: "",
          password: ""
        };
      }
    
      onChangeEmail(e) {
        this.setState({
          email: e.target.value,
        });
      }
    
      onChangePassword(e) {
        this.setState({
          password: e.target.value,
        });
      }
    
      handleLogin(e) {
        e.preventDefault();
        
        this.form.validateAll();
        const { dispatch} = this.props;
    
        if (this.checkBtn.context._errors.length === 0) {
          dispatch(login(this.state.email, this.state.password))
            .then(() => {
              window.location.reload();
            })
        }
        
      }

    render(){
        const { isLoggedIn, message } = this.props;

        if (isLoggedIn) {
        return <Redirect to="/profile"/>;
        }
        return (
            <Form className="mt-8" onSubmit={this.handleLogin}
                ref={(c) => {
                this.form = c;
                }
            }>
                <input type="hidden" name="remember" value="true" />
                <div className="rounded shadow-sm">
                    <div>
                        <Input aria-label="Email address" value={this.state.email} onChange={this.onChangeEmail} validations={[required]} name="email"  type="email" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5" placeholder="Email address"/>
                    </div>
                    <div className="mt-4">
                        <Input aria-label="Password" value={this.state.password} onChange={this.onChangePassword} validations={[required]} name="password"  type="password" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5" placeholder="Password" />
                    </div>
                </div>

                <div className="mt-2 flex flex-row justify-between" >
                    <div className="flex  items-center">
                        <input  name="remember-me" type={"checkbox"}/><span className="text-xs text-purple-900">Remember me</span>
                    </div>
                    <div className="flex items-center">
                        <a href="#" className="font-medium text-purple-900 text-xs hover:text-purple-700 focus:outline-none focus:underline transition ease-in-out duration-150">
                            Forgot your password?
                        </a>
                    </div>
                </div>
                <div className="mt-3 pb-3">
                    <button type="submit" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                        Sign in
                    </button>
                </div> 
                {message && (
                    <div className="form-group">
                        <div className="alert alert-danger" role="alert">
                        {message}
                        </div>
                    </div>
                )}
                <CheckButton
                    style={{ display: "none" }}
                    ref={(c) => {
                        this.checkBtn = c;
                    }}
                />
                
            </Form>
        )
    }
}

function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    const { message } = state.message;
    return {
      isLoggedIn,
      message
    };
  }
  
  export default connect(mapStateToProps)(LoginForm);