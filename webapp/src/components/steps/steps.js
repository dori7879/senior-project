import React from 'react';

const Steps = () =>{
    return(
        <div>
            <div className="flex justify-center text-purple-900 font-bold text-xl pt-4  border-t-2 border-purple-200">
                How to use EasySubmit?
            </div>
            <div className="flex justify-center">
                <section className="text-gray-600 body-font">
                    <div className="container px-5 py-4 mx-auto flex flex-wrap">
                        <div className="flex relative pt-10 pb-20 sm:items-center md:w-2/3 mx-auto">
                        <div className="h-full w-6 absolute inset-0 flex items-center justify-center">
                            <div className="h-full w-1 bg-purple-200 pointer-events-none"></div>
                        </div>
                        <div className="flex-shrink-0 w-6 h-6 rounded-full mt-8 sm:mt-0 inline-flex items-center justify-center bg-purple-600 text-purple-100 relative z-10 title-font font-medium text-sm">1</div>
                        <div className="flex-grow md:pl-8 pl-6 flex sm:items-center items-start flex-col sm:flex-row">
                            <div className="flex-grow sm:pl-6 mt-6 sm:mt-0">
                            <h2 className="font-medium title-font text-purple-700 mb-1 text-xl">Choose type of the assignment</h2>
                            <p className="leading-relaxed">Create a homework (quiz or attendance) page by clicking on the corresponding button on the Home page</p>
                            </div>
                        </div>
                        </div>
                        <div className="flex relative pb-20 sm:items-center md:w-2/3 mx-auto">
                        <div className="h-full w-6 absolute inset-0 flex items-center justify-center">
                            <div className="h-full w-1 bg-purple-200 pointer-events-none"></div>
                        </div>
                        <div className="flex-shrink-0 w-6 h-6 rounded-full mt-10 sm:mt-0 inline-flex items-center justify-center bg-purple-600 text-purple-100 relative z-10 title-font font-medium text-sm">2</div>
                        <div className="flex-grow md:pl-8 pl-6 flex sm:items-center items-start flex-col sm:flex-row">

                            <div className="flex-grow sm:pl-6 mt-6 sm:mt-0">
                            <h2 className="font-medium title-font text-purple-700 mb-1 text-xl">Fill in form</h2>
                            <p className="leading-relaxed">Enter the title and the description of the assignment. Attach files if needed.</p>
                            </div>
                        </div>
                        </div>
                        <div className="flex relative pb-20 sm:items-center md:w-2/3 mx-auto">
                        <div className="h-full w-6 absolute inset-0 flex items-center justify-center">
                            <div className="h-full w-1 bg-purple-200 pointer-events-none"></div>
                        </div>
                        <div className="flex-shrink-0 w-6 h-6 rounded-full mt-10 sm:mt-0 inline-flex items-center justify-center bg-purple-600 text-purple-100 relative z-10 title-font font-medium text-sm">3</div>
                        <div className="flex-grow md:pl-8 pl-6 flex sm:items-center items-start flex-col sm:flex-row">
                            
                            <div className="flex-grow sm:pl-6 mt-6 sm:mt-0">
                            <h2 className="font-medium title-font text-purple-700 mb-1 text-xl">Generate Link</h2>
                            <p className="leading-relaxed">Click on the button to generate a link for the homework (quiz or attendance) page and new page will be created.</p>
                            </div>
                        </div>
                        </div>
                        <div className="flex relative pb-10 sm:items-center md:w-2/3 mx-auto">
                        <div className="h-full w-6 absolute inset-0 flex items-center justify-center">
                            <div className="h-full w-1 bg-purple-200 pointer-events-none"></div>
                        </div>
                        <div className="flex-shrink-0 w-6 h-6 rounded-full mt-10 sm:mt-0 inline-flex items-center justify-center bg-purple-600 text-purple-100 relative z-10 title-font font-medium text-sm">4</div>
                        <div className="flex-grow md:pl-8 pl-6 flex sm:items-center items-start flex-col sm:flex-row">
                            
                            <div className="flex-grow sm:pl-6 mt-6 sm:mt-0">
                            <h2 className="font-medium title-font text-purple-700 mb-1 text-xl">Share Link</h2>
                            <p className="leading-relaxed">Copy the generated link and share it with your students. Track their submissions.</p>
                            </div>
                        </div>
                        </div>
                    </div>
                </section>
            </div>
        </div>
    )
}

export default Steps;