import Footer from '../footer';
import Header from '../header';
import React from 'react';
import Steps from '../steps';

const HomePage = () => {
    return(
        <div>
            <Header />
            <div className="min-h-screen flex items flex-col justify-center bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                <div className="flex items-start justify-center pt-6"> 
                    <div className="flex justify-around pb-12">
                        <div className=" mr-10">
                            <img src="https://img.pngio.com/online-education-png-education-png-456_413.png" alt="" className=""/>
                        </div>
                        <div className="lg:mt-20 ml-10 sm:mt-10">
                                <button type="button" className="h-10 items-center mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-600 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                    Create homework page
                                </button>
                                <button type="button" className="h-10 items-center mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-600 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                    Create quiz page
                                </button>
                                <button type="button" className="h-10 items-center mb-2 relative w-full flex justify-center pt-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-600 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                    Create attendance page
                                </button>
                        </div>
                    </div>
                </div>
                <div className="flex justify-center flex-col items-center border-t-2 border-purple-200 pb-8">
                    <div className="text-purple-900 font-bold text-xl pt-4 pb-2">Introducing EasySubmit</div>
                    <div className="border-2 border-purple-400 rounded-lg w-2/3 bg-purple-200">
                        <p className="flex text-gray-700 justify-center  mx-4  my-2">A progressive web application (PWA) that allows teachers/instructors to post homework (quiz or attendance) simply by opening the website and generating a link for the new homework page. This link could then be shared with students to submit their solutions.</p>
                    </div>
                </div>
                <Steps />
            </div>
            <Footer />
        </div>
    )
}

export default HomePage;