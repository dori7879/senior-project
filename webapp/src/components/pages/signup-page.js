import Footer from '../footer';
import Header from '../header';
import React from 'react';
import RegistrationForm from '../registration-form/registration-form';

class SignUpPage extends React.Component{
   
    render(){
        return(
            <div>
            <Header />
            <div className="min-h-screen flex items-center justify-center bg-purple-100 py-8 px-4 sm:px-6 lg:px-8">
                <div className="border border-purple-300">
                    <div className="max-w-md w-64 mx-4">
                        <div>
                            <h2 className="mt-6 text-center text-2xl uppercase leading-7 font-extrabold text-purple-900">
                                Register
                            </h2>
                        </div>
                        <RegistrationForm />
                    </div>    
                </div>
            </div>
            <Footer />
         </div> 
        )
           

    }
}

export default SignUpPage;