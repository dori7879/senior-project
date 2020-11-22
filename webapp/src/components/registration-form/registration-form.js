import CheckButton from "react-validation/build/button";
import Form from "react-validation/build/form";
import Input from "react-validation/build/input";
import React from "react";
import { connect } from "react-redux";
import { isEmail } from "validator";
import { register } from "../../actions/auth";

const required = (value) => {
    if (!value) {
        return (
        <div className="border border-rounded bg-red-200 text-red-800 text-center" role="alert">
            This field is required!
        </div>
        );
    }
};
  
const email = (value) => {
    if (!isEmail(value)) {
        return (
        <div className="border border-rounded bg-red-200 text-red-800 text-center" role="alert">
            This is not a valid email.
        </div>
        );
    }
};
  
const vfirstname = (value) => {
    if (value.length < 2 || value.length > 20) {
        return (
        <div className="border border-rounded bg-red-200 text-red-800 text-center" role="alert">
            The username must be between 2 and 20 characters.
        </div>
        );
    }
};
  
const vlastname = (value) => {
    if (value.length < 2 || value.length > 20) {
        return (
        <div className="border border-rounded bg-red-200 text-red-800 text-center" role="alert">
            The username must be between 2 and 20 characters.
        </div>
        );
    }
};
const vpassword = (value) => {
    if (value.length < 6 || value.length > 40) {
        return (
        <div className="border border-rounded bg-red-200 text-red-800 text-center" role="alert">
            The password must be between 6 and 40 characters.
        </div>
        );
    }
};
  

class RegistrationForm extends React.Component {
    constructor(props) {
        super(props);
        this.handleRegister = this.handleRegister.bind(this);
        this.onChangeEmail = this.onChangeEmail.bind(this);
        this.onChangeFirstName = this.onChangeFirstName.bind(this);
        this.onChangeLastName = this.onChangeLastName.bind(this);
        this.onChangePassword = this.onChangePassword.bind(this);
    
    this.state = {
            firstName: "",
            lastName: "",
            email: "",
            password: "",
            successful: false,
        };
    }

    onChangeFirstName(e) {
        this.setState({
            firstName: e.target.value,
        });
    }

    onChangeLastName(e) {
        this.setState({
            lastName: e.target.value,
        });
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

    handleRegister(e) {
        e.preventDefault();

        this.setState({
            successful: false,
        });

        this.form.validateAll();

        if (this.checkBtn.context._errors.length === 0) {
            this.props
            .dispatch(
                register(this.state.firstName, this.state.lastName, this.state.email, this.state.password)
            )
            .then(() => {
                this.setState({
                successful: true,
                });
            })
            .catch(() => {
                this.setState({
                successful: false,
                });
            });
        }
    }
    render() {
        const { message } = this.props;
        return(
            <Form className="mt-8" onSubmit={this.handleRegister}
                ref={(c) => {
                this.form = c;
                }}
             >
                {!this.state.successful && (
                    <div>
                        <input type="hidden" name="remember" value="true" />
                        <div className="rounded shadow-sm">
                            <div>
                                <label className="text-sm text-gray-500">First Name</label>
                                <Input aria-label="First name"  value={this.state.firstName} onChange={this.onChangeFirstName} validations={[required, vfirstname]} name="firstname" type="firstname" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5" />
                            </div>
                            <div >
                                <label className="text-sm text-gray-500">Last Name</label>
                                <Input aria-label="Last Name" value={this.state.lastName} onChange={this.onChangeLastName} validations={[required, vlastname]} name="lastname" type="lastname" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5"/>
                            </div>
                            <div>
                                <label className="text-sm text-gray-500">Email address</label>
                                <Input aria-label="Email address" value={this.state.email} onChange={this.onChangeEmail} validations={[required, email]} name="email" type="email" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5"/>
                            </div>
                            <div>
                                <label className="text-sm text-gray-500">Password</label>
                                <Input aria-label="Password" value={this.state.password} onChange={this.onChangePassword} validations={[required, vpassword]}name="password" type="password" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5"/>
                            </div>
                        </div>
                        <div className="mt-4 pb-4">
                            <button type="submit" className="mb-2 relative w-full flex justify-center py-2 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                Sign up
                            </button>
                        </div>
                        </div>
                )}
                
                {message && (
                    <div className="form-group">
                        <div className={ this.state.successful ? "alert alert-success" : "alert alert-danger" } role="alert">
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
        );
    }
}

function mapStateToProps(state) {
    const { message } = state.message;
    return {
      message,
    };
  }
  
  export default connect(mapStateToProps)(RegistrationForm);