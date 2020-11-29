import CKEditor from 'ckeditor4-react';
import Footer from '../footer';
import Header from '../header';
import React from 'react';
import axios from 'axios';
import { submitHomework } from "../../actions/homework";

class StudentHwPage extends React.Component{
    
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeFullName = this.onChangeFullName.bind(this);
        this.onChangeAnswer = this.onChangeAnswer.bind(this);
        this.onChangeAttachments = this.onChangeAttachments.bind(this);

        this.state = {
            fullName: "",
            answer: "",
            attachments: [],
            submitDate: new Date(),
            successful: false,
            course_title: "",
            title: "",
            description: "",
            files: [],
            closeDate: new Date(),
            grade: "",
            comments: ""
        };          
    }
    onChangeFullName(e) {
        this.setState({
          fullName: e.target.value,
        });
    }   
    onChangeAnswer(e) {
        this.setState({
        answer: e.editor.getData()
        });
    }
    onChangeAttachments(e) {
        this.setState({
          files: e.target.value[0],
         });
    }
    handleSubmit(e){
        e.preventDefault();
        const { randomStr } = this.props.match.params;
        this.setState({
            isClicked: true,
            submitDate:new Date()
        });
        const { dispatch} = this.props;
        dispatch(submitHomework( this.state.fullName, this.state.answer, this.state.submitDate, this.state.grade, this.state.comments))
            .then(() => {
                this.setState({
                    isClicked: true,
                    successful: true,
                });
            })
            .catch(() => {
                this.setState({
                    successful: false,
                });
            });
    }
    componentDidMount () {
        const { randomStr } = this.props.match.params;
        axios.get(`/api/v1/homework-page/student/${randomStr}`)
            .then((response) => {
                if (response.data) {
                    this.setState({
                        course_title: response.data.course_title,
                        title: response.data.title,
                        description: response.data.content,
                        closeDate: response.data.closed_at
                    })
                }
            })
        /*const { dispatch} = this.props;
        dispatch(fetchHomework(randomStr))
            .then((response) => {
                if (response.data) {
                    this.setState({
                        course_title: response.data.course_title,
                        title: response.data.title,
                        description: response.data.content,
                        closeDate: response.data.closed_at
                    })
                }
            
        })*/
    }
    render(){
        const isEmptyDesc = this.state.description.trim() === "";
        const isEmptyFile = this.state.files.length === 0; 
        return(
            <div>
                <Header />
                <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                    <div className="w-3/4 flex justify-center flex-col items-center pb-4">
                        <div className="text-purple-900 font-bold text-xl pt-2">Homework</div>
                            <div className="w-full  flex flex-row items items-center border border-purple-300 rounded-t ">
                                <div className="w-full border border-purple-300 bg-purple-300 rounded-tl p-4">
                                    <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-3 px-4 pt-1">
                                        <strong>Course Title:</strong> <span className="text-purple-900">{this.state.course_title}</span></h2>  
                                    <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                        <strong>Homework Title:</strong> <span className="text-purple-900">{this.state.title}</span></h2>   
                                </div>
                                <div>
                                <div className="flex flex-row items">
                                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-4 px-4 pt-1" >
                                        Name:
                                    </label>
                                    <input onChange={this.onChangeFullName} className="mr-4 mb-3 text-gray-700 border border-purple-400  rounded  text-xs px-1 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="Enter your full name" />
                                </div>
                                <div className="flex flex-row items items-center">
                                    <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                    <strong>Time Remaining:</strong> <span className="text-purple-900">15min</span></h2>
                                </div>  
                                
                            </div>     
                        </div>
                        <div  className="w-full border border-purple-300 rounded-b bg-purple-300 p-4">
                            {
                                isEmptyDesc ? null :
                                <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                    <strong>Description:</strong><br></br> <span className="text-purple-900">{this.state.description}</span></h2>
                            }
                            {
                                isEmptyFile ? null :
                                <h1> Attachments</h1> 
                            }
                            <form className=" onSubmit={this.handleSubmit}">
                                <div className="flex flex-col pb-2 mx-4">
                                    <h1 className="block uppercase tracking-wide text-purple-900 text-xs font-bold mb-2 mt-4 px-4 pt-1 text-center" >
                                        Answer
                                    </h1>
                                    <CKEditor
                                        data={this.state.answer}
                                        onChange={this.onChangeAnswer}
                                    />
                                </div>
                                <div className="flex flex-row ml-2 pb-2 items items-center">
                                   <label className="block uppercase tracking-wide  text-gray-700 text-xs font-bold mb-2 px-2 pt-1" >
                                        Attach files
                                    </label>
                                    <input onChange={this.onChangeAttachments} className="border border-purple-400 text-xs text-gray-700 w-1/2 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="files" type="file" multiple />
                                </div>
                                <div className="flex justify-center">
                                    <button type="submit" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                         Submit
                                    </button>
                                </div>
                            </form>
                        </div>
                         
                    </div>
                </div>   
                <Footer />
            </div>
        )
    }
    }



export default StudentHwPage;