import Footer from '../footer';
import Header from '../header';
import React from 'react';
import { Redirect } from 'react-router-dom';
import { useSelector } from "react-redux";

const MyProfilePage = () => {
    
    const { access_token, first_name, last_name, email, role } = useSelector((state) => state.auth);

    if (!access_token) {
      return <Redirect to="/signin"/>;
    }    

    return(
        <div>
            <Header />
            <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">  
                <div className="flex justify-center flex-col items-center pb-4">
                    <div className="text-purple-900 font-bold text-xl pt-4">My profile</div>
                </div>        
                    <div className="w-1-2 border border-purple-300 rounded bg-purple-300 p-4">
                         <p className="text-purple-900">
                            <strong>First Name: </strong> {first_name}
                        </p>
                        <p className="text-purple-900">
                            <strong>Last Name: </strong> {last_name}
                        </p>
                        <p  className="text-purple-900">
                            <strong>Email:</strong> {email}
                        </p>
                        <p  className="text-purple-900">
                            <strong>Role:</strong> {role}
                        </p>
                        <p className="text-purple-900 mt-3">
                            <strong >Change Password </strong> 
                        </p>
                       
                        <form className=""> 
                            <div className="flex flex-row items items-center pb-2">
                                <div className="flex flex-col">
                                    <label className="w-full block uppercase tracking-wide text-gray-700 text-xs font-bold ml-5 mb-2 px-4 pt-1" >
                                        New Password
                                    </label>
                                    <input className=" text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 mr-6 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
                                </div>
                                <div className="flex flex-col">
                                    <label className="w-full block uppercase tracking-wide text-gray-700 text-xs font-bold ml-3 mb-2 px-4 pt-1" >
                                        Confirm Password
                                    </label>
                                    <input className="text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 mx-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
                                </div>
                            </div> 
                            <button type="submit" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                Save
                            </button>
                        </form>                         
                        
                        <p  className="text-gray-700 mt-3">
                            <strong className="text-purple-900"> My Classes </strong> 
                        </p>
                        <button type="submit" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                10A
                        </button>
                        <button type="submit" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                11B
                        </button>
                    </div>
                   
            </div>   
            <Footer />
        </div>
    )
}

export default MyProfilePage;
