import Footer from '../footer';
import Header from '../header';
import React from 'react';
import { Redirect } from 'react-router-dom';
import { useSelector } from "react-redux";

const LinkPage = () => {
    const { homework } = useSelector((state) => state.homework);

    if (!homework) {
      return <Redirect to="/homework"/>;
    }      

    return(
        <div>
            <Header />
            <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                <div className="flex justify-center flex-col items-center pb-4">
                    <div className="text-purple-900 font-bold text-xl pt-4">Links</div>
                    <div className="container">
                        <p>
                            <strong>Link for Students:</strong> {homework.studentLink}
                        </p>
                        <p>
                            <strong>Link for Teacher:</strong> {homework.teacherLink}
                        </p>
                    </div>
                </div>
            </div>   
            <Footer />
        </div>
    )
}

export default LinkPage;