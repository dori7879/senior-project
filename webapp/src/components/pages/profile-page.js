import Footer from '../footer';
import Header from '../header';
import React from 'react';
import { Redirect } from 'react-router-dom';
import { useSelector } from "react-redux";

const MyProfilePage = () => {
    const { user: currentUser } = useSelector((state) => state.auth);

    if (!currentUser) {
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
                            <strong>Token:</strong> {currentUser.token.substring(0, 20)} ...{" "}
                            {currentUser.token.substr(currentUser.token.length - 20)}
                        </p>
                        <p>
                            <strong>Email:</strong> {currentUser.email}
                        </p>
                    </div>
                </div>
            </div>   
            <Footer />
        </div>
    )
}

export default MyProfilePage;