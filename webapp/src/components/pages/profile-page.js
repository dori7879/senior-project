import Footer from '../footer';
import Header from '../header';
import React from 'react';
import { Redirect } from 'react-router-dom';
import { useSelector } from "react-redux";

const MyProfilePage = () => {
    const { access_token } = useSelector((state) => state.auth);

    if (!access_token) {
      return <Redirect to="/signin"/>;
    }      

    return(
        <div>
            <Header />
            <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                <div className="flex justify-center flex-col items-center pb-4">
                    <div className="text-purple-900 font-bold text-xl pt-4">My profile</div>
                    <div className="container">
                        <p>
                            <strong>Token:</strong> {access_token}
                        </p>
                        <p>
                            <strong>Email:</strong> {access_token}
                        </p>
                    </div>
                </div>
            </div>   
            <Footer />
        </div>
    )
}

export default MyProfilePage;