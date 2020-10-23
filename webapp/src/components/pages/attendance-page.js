import Footer from '../footer';
import Header from '../header';
import React from 'react';

const AttendancePage = () => {
    return(
        <div>
            <Header />
            <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                <div className="flex justify-center flex-col items-center pb-4">
                    <div className="text-purple-900 font-bold text-xl pt-4 ">Creating an Attendance check</div>
                </div>
                <form className="w-3/4 border border-purple-300 rounded bg-purple-300 p-4" >
                    <div className="flex flex-row items items-center pb-2">
                        <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                            Title*
                        </label>
                        <input className="text-gray-700 border border-purple-400 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="Enter the title" />
                    </div>
                    <div className="flex flex-row pb-2 items items-center">
                        <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                            Open date and time*
                        </label>
                        <input className="block border border-purple-400 rounded text-gray-700 w-1/4 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="open-date" name="trip-start"  type="date" value="2020-10-16" min="2020-10-16" max="2021-12-31" />
                        <input className="block border border-purple-400 rounded text-gray-700 w-1/6 rounded py-1 px-2 ml-3 leading-tight focus:outline-none focus:bg-white" id="open-time" name="trip-start" type="time" min="00:00" max="11:59" />
                    </div>
                    <div className="flex flex-row pb-2 items items-center">
                        <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                            Close date and time*
                        </label>
                        <input className="block border border-purple-400 rounded text-gray-700 w-1/4 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="close-date" name="trip-start"  type="date" value="2020-10-16" min="2020-10-16" max="2021-12-31" />
                        <input className="block border border-purple-400 rounded text-gray-700 w-1/6 rounded py-1 px-2 ml-3 leading-tight focus:outline-none focus:bg-white" id="close-time" name="trip-start" type="time" min="00:00" max="11:59" />
                    </div>
                    <div className="flex flex-col pb-2 border border-purple-400 rounded mb-2 w-1/2">
                        <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                            Choose an attendance type*
                        </label>
                        <div>
                            <div className="flex flex-row pl-4">
                                <input type="radio" id="choice1" name="attendance-type" value="Face-recognition"
                                        checked />
                                <label for="choice1" className="text-gray-700 rounded py-1 px-2 leading-tight">Face recognition</label>
                            </div>
                            <div className="flex flex-row pl-4">
                                <input type="radio" id="choice2" name="attendance-type" value="Generated-ID"
                                        />
                                <label for="choice2" className="text-gray-700 rounded py-1 px-2 leading-tight">Generated ID</label>
                            </div>
                            <div className="flex flex-row pl-4">
                                <input type="radio" id="choice3" name="attendance-type" value="QR-code"
                                        />
                                <label for="choice3" className="text-gray-700 rounded py-1 px-2 leading-tight">QR code</label>
                            </div>
                        </div>
                    </div>
                    <div className="flex justify-center">
                        <input className="shadow bg-purple-800 hover:bg-purple-500 focus:shadow-outline focus:outline-none text-white font-bold py-2 px-4 rounded" type="submit" value="Generate a link"/>
                    </div>
                    <div className="flex flex-col items items-center justify-center mt-2">
                          <textarea rows="1" className="text-gray-700 border border-purple-300 w-1/2 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="description" type="text" />  
                    </div>
                </form>
                <small className="mt-4">* - Required field </small>  
            </div>
            <Footer />
        </div>
    )
}

export default AttendancePage;