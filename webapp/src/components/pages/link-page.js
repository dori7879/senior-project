import Footer from '../footer';
import Header from '../header';
import React from 'react';
import { connect } from "react-redux";

class LinkPage extends React.Component {
       

    constructor(props) {
        super(props);
    
        this.state = {
             copySuccess_1: '', 
             copySuccess: ''
        }
      }
    
      copyToClipboardStudent = (e) => {
        this.textArea1.select();
        document.execCommand('copy');
        e.target.focus();
        this.setState({ copySuccess_1: 'Copied!' });
      };
      copyToClipboardTeacher = (e) => {
        this.textArea.select();
        document.execCommand('copy');
        e.target.focus();
        this.setState({ copySuccess: 'Copied!' });
      };

    render(){
        const {teacher_link, student_link} = this.props;
        const url_teacher = `https://radiant-inlet-12251.herokuapp.com/teacher-hw-page/${teacher_link}`;
        const url_student = `https://radiant-inlet-12251.herokuapp.com/student-hw-page/${student_link}`;
        
        return(
            <div>
                <Header />
                <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                    <div className="flex justify-center flex-col items-center pb-4">
                        <div className="text-purple-900 font-bold text-xl pt-4">Links</div>
                        
                        <div className="flex flex-row ">
                            <div className="mr-20">
                                 <h1 className="text-gray-700 font-bold text-xl pt-4">For students:</h1>
                                {
                                document.queryCommandSupported('copy') &&
                                <div className="text-purple-900 font-bold text-xl pt-4">
                                    <button className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out" onClick={this.copyToClipboardStudent}>Copy</button> 
                                    {this.state.copySuccess_1}
                                </div>
                                }
                                <form>
                                    <textarea
                                        ref={(textarea1) => this.textArea1 = textarea1}
                                        defaultValue={url_student}
                                        rows="10"
                                    />
                                </form>
                            </div>
                            <div >
                                <h1 className="text-gray-700 font-bold text-xl pt-4">For teacher:</h1>
                                {
                                document.queryCommandSupported('copy') &&
                                <div className="text-purple-900 font-bold text-xl pt-4">
                                    <button className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out" onClick={this.copyToClipboardTeacher}>Copy</button> 
                                    {this.state.copySuccess}
                                </div>
                                }
                                <form>
                                    <textarea
                                        ref={(textarea) => this.textArea = textarea}
                                        defaultValue={url_teacher}
                                        rows="10"
                                    />
                                </form>
                            </div>
                        </div>
                       
                    </div>
                </div>   
                <Footer />
            </div>
        )
    }
}

function mapStateToProps(state) {
    const { teacher_link, student_link} = state.homework;
    return {
      teacher_link, student_link
    };
  }

export default connect(mapStateToProps)(LinkPage);