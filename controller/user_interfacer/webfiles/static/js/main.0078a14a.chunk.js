(this["webpackJsonpencodarr-react-frontend"]=this["webpackJsonpencodarr-react-frontend"]||[]).push([[0],{36:function(e,t,a){},62:function(e,t,a){},81:function(e,t,a){},82:function(e,t,a){},86:function(e,t,a){},88:function(e,t,a){},89:function(e,t,a){},90:function(e,t,a){},91:function(e,t,a){"use strict";a.r(t);a(57);var s=a(1),i=a.n(s),n=a(27),r=a.n(n),c=(a(62),a(11)),o=a(12),l=a(14),d=a(13),h=a(31),j=a(22),b=a(10),u=a.n(b),p=a(19),x=a(29),O=a(24),v=a(8),m=a(56),f=a(46),g=(a(81),a(36),a.p+"static/media/Info-I.ffc9d3a2.svg"),k=(a(82),a(0));function _(e){return Object(k.jsx)("img",{className:"queue-icon",src:e.location,alt:e.alt,height:"20px",title:e.title})}var y=a.p+"static/media/terminalIcon.5147de0e.svg";function C(e){return Object(k.jsx)(_,{location:y,alt:"Terminal",title:e.title})}var w=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).timerID=void 0,s.state={jobs:[],waitingOnServer:!0,showModal:!1,waitingRunnersText:""},s.timerID=setTimeout((function(){}),Number.POSITIVE_INFINITY),clearInterval(s.timerID),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;this.tick(),this.timerID=setInterval((function(){return e.tick()}),2e3)}},{key:"componentWillUnmount",value:function(){clearInterval(this.timerID)}},{key:"tick",value:function(){var e=this;u.a.get("/api/web/v1/running").then((function(t){var a=t.data.jobs;void 0!==a?(a.sort((function(e,t){return parseFloat(e.status.percentage)>parseFloat(t.status.percentage)?-1:1})),e.setState({jobs:a,waitingOnServer:!1})):console.error("Response from /api/web/v1/running returned undefined for data.jobs")})).catch((function(e){console.error("Request to /api/web/v1/running failed with error: ".concat(e))})),u.a.get("/api/web/v1/waitingrunners").then((function(t){if(0===t.data.Runners.length)e.setState({waitingRunnersText:"No waiting runners"});else{var a=t.data.Runners.toString();1!==t.data.Runners.length&&(a=a.slice(1)),e.setState({waitingRunnersText:a})}})).catch((function(e){console.error("Request to /api/web/v1/waitingrunners failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this,t=function(){return e.setState({showModal:!1})},a=this.state.jobs.map((function(e){return Object(k.jsx)(S,{fps:e.status.fps,uuid:e.job.uuid,filename:e.job.path,progress:e.status.percentage,runnerName:e.runner_name,stageValue:e.status.stage,jobElapsedTime:e.status.job_elapsed_time,stageElapsedTime:e.status.stage_elapsed_time,stageEstimatedTimeRemaining:e.status.stage_estimated_time_remaining,command:e.job.command.join(" ")},e.job.uuid)}));return Object(k.jsxs)("div",{children:[Object(k.jsx)("img",{className:"info-i",src:g,alt:"",height:"20px",onClick:function(){return e.setState({showModal:!0})}}),0!==a.length?a:Object(k.jsx)("h5",{className:"text-center",children:"No running jobs"}),Object(k.jsxs)(v.a,{show:this.state.showModal,onHide:t,children:[Object(k.jsx)(v.a.Header,{closeButton:!0,children:Object(k.jsx)(v.a.Title,{children:"Waiting Runners"})}),Object(k.jsx)(v.a.Body,{children:this.state.waitingRunnersText}),Object(k.jsx)(v.a.Footer,{children:Object(k.jsx)(p.a,{variant:"secondary",onClick:t,children:"Close"})})]})]})}}]),a}(i.a.Component);function S(e){return Object(k.jsxs)("div",{children:[Object(k.jsxs)(x.a,{style:{padding:"1rem"},children:[Object(k.jsxs)(x.a.Header,{className:"text-center",children:[Object(k.jsxs)("div",{className:"file-image-container",children:[Object(k.jsx)("h5",{children:e.filename}),Object(k.jsx)(C,{title:e.command})]}),Object(k.jsxs)("h6",{children:["Stage: ",e.stageValue]}),Object(k.jsxs)("h6",{children:["Runner: ",e.runnerName]})]}),Object(k.jsx)(m.a,{className:"progress-bar-style",animated:!0,now:parseFloat(e.progress),label:"".concat(e.progress,"%")}),Object(k.jsxs)(f.a,{children:[Object(k.jsx)(O.a,{children:Object(k.jsx)("h6",{className:"text-right",children:"Job Elapsed Time:"})}),Object(k.jsx)(O.a,{children:Object(k.jsx)("p",{children:e.jobElapsedTime})}),Object(k.jsx)(O.a,{children:Object(k.jsx)("h6",{className:"text-right",children:"FPS:"})}),Object(k.jsx)(O.a,{children:Object(k.jsx)("p",{children:e.fps})})]}),Object(k.jsxs)(f.a,{children:[Object(k.jsx)(O.a,{children:Object(k.jsx)("h6",{className:"text-right",children:"Stage Elapsed Time:"})}),Object(k.jsx)(O.a,{children:Object(k.jsx)("p",{children:e.stageElapsedTime})}),Object(k.jsx)(O.a,{children:Object(k.jsx)("h6",{className:"text-right",children:"Stage Estimated Time Remaining:"})}),Object(k.jsx)(O.a,{children:Object(k.jsx)("p",{children:e.stageEstimatedTimeRemaining})})]})]}),Object(k.jsx)("div",{className:"smol-spacer"})]})}var N=a(28),I=a(15),T=a(4),H=a(39),E=(a(86),a.p+"static/media/addLibraryIcon.dd5f1d29.svg"),R=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).timerID=void 0,s.state={libraries:[],waitingOnServer:!0,showCreateLibModal:!1},s.timerID=setTimeout((function(){}),Number.POSITIVE_INFINITY),clearInterval(s.timerID),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;this.tick(),this.timerID=setInterval((function(){return e.tick()}),2e3)}},{key:"componentWillUnmount",value:function(){clearInterval(this.timerID)}},{key:"tick",value:function(){var e=this;u.a.get("/api/web/v1/libraries").then((function(t){200===t.status&&e.setState({libraries:t.data.IDs})})).catch((function(e){console.error("Request to /api/web/v1/libraries failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this,t=this.state.libraries.map((function(e){return Object(k.jsxs)("div",{children:[Object(k.jsx)(V,{id:e}),Object(k.jsx)("div",{className:"smol-spacer"})]},e)}));return Object(k.jsxs)(k.Fragment,{children:[Object(k.jsx)("img",{className:"add-lib-ico",src:E,alt:"",height:"20px",onClick:function(){e.setState({showCreateLibModal:!0})}}),Object(k.jsx)(P,{show:this.state.showCreateLibModal,closeHandler:function(){e.setState({showCreateLibModal:!1})}}),Object(k.jsx)("div",{className:"smol-spacer"}),t]})}}]),a}(i.a.Component),V=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={folder:"",priority:"",fs_check_interval:"",path_masks:"",queue:[],target_video_codec:"HEVC",create_stereo_audio:!0,skip_hdr:!0,showEditModal:!1,showQueueModal:!1},s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){this.getLibraryData()}},{key:"getLibraryData",value:function(){var e=this;u.a.get("/api/web/v1/library/".concat(this.props.id)).then((function(t){var a=JSON.parse(t.data.command_decider_settings);e.setState({folder:t.data.folder,priority:t.data.priority,fs_check_interval:t.data.fs_check_interval,path_masks:t.data.path_masks.join(","),queue:t.data.queue.Items,target_video_codec:a.target_video_codec,create_stereo_audio:a.create_stereo_audio,skip_hdr:a.skip_hdr})})).catch((function(t){console.error("Request to /api/web/v1/library/".concat(e.props.id," failed with error: ").concat(t))}))}},{key:"render",value:function(){var e=this;return Object(k.jsxs)(k.Fragment,{children:[Object(k.jsxs)(x.a,{children:[Object(k.jsx)(x.a.Header,{className:"text-center",children:Object(k.jsx)("h5",{children:this.state.folder})}),Object(k.jsxs)("p",{className:"text-center",children:["Priority: ",this.state.priority]}),Object(k.jsxs)("p",{className:"text-center",children:["File System Check Interval: ",this.state.fs_check_interval]}),Object(k.jsxs)("p",{className:"text-center",children:["Target Video Codec: ",this.state.target_video_codec]}),Object(k.jsxs)("p",{className:"text-center",children:["Create Stereo Audio Track: ",this.state.create_stereo_audio?"True":"False"]}),Object(k.jsxs)("p",{className:"text-center",children:["Skip HDR Files: ",this.state.skip_hdr?"True":"False"]}),0!==this.state.path_masks.length?Object(k.jsxs)("p",{className:"text-center",children:["Path Masks: ",this.state.path_masks]}):null,Object(k.jsx)(p.a,{variant:"secondary",onClick:function(){e.setState({showQueueModal:!0})},children:"Queue"}),Object(k.jsx)(p.a,{variant:"primary",onClick:function(){e.setState({showEditModal:!0})},children:"Edit"})]}),this.state.showEditModal?Object(k.jsx)(L,{show:!0,closeHandler:function(){e.setState({showEditModal:!1}),e.getLibraryData()},id:this.props.id,folder:this.state.folder,priority:this.state.priority,fs_check_interval:this.state.fs_check_interval,path_masks:this.state.path_masks,target_video_codec:this.state.target_video_codec,create_stereo_audio:this.state.create_stereo_audio,skip_hdr:this.state.skip_hdr}):null,this.state.showQueueModal?Object(k.jsx)(D,{show:!0,closeHandler:function(){e.setState({showQueueModal:!1})},queue:this.state.queue}):null]})}}]),a}(i.a.Component),P=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={folder:"",priority:"",fs_check_interval:"",path_masks:"",target_video_codec:"HEVC",create_stereo_audio:!0,skip_hdr:!0},s.submitLib=s.submitLib.bind(Object(N.a)(s)),s}return Object(o.a)(a,[{key:"submitLib",value:function(){var e=this,t={folder:this.state.folder,priority:parseInt(this.state.priority),fs_check_interval:this.state.fs_check_interval,path_masks:this.state.path_masks.split(","),pipeline:{target_video_codec:this.state.target_video_codec,create_stereo_audio:this.state.create_stereo_audio,skip_hdr:this.state.skip_hdr}};u.a.post("/api/web/v1/library/new",t).then((function(){e.props.closeHandler()})).catch((function(e){console.error("/api/web/v1/library/new failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this;return Object(k.jsx)("div",{children:Object(k.jsxs)(v.a,{show:this.props.show,onHide:this.props.closeHandler,children:[Object(k.jsx)(v.a.Header,{closeButton:!0,children:Object(k.jsx)(v.a.Title,{children:"Create New Library"})}),Object(k.jsxs)(v.a.Body,{children:[Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Folder"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"/home/user/lib1","aria-label":"folder","aria-describedby":"basic-addon1",onChange:function(t){e.setState({folder:t.target.value})},value:this.state.folder})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Priority"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"0","aria-label":"priority","aria-describedby":"basic-addon1",onChange:function(t){e.setState({priority:t.target.value})},value:this.state.priority})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"File System Check Interval"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"fs_check_interval","aria-describedby":"basic-addon1",onChange:function(t){e.setState({fs_check_interval:t.target.value})},value:this.state.fs_check_interval})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Target Video Codec"})}),Object(k.jsxs)(I.a,{className:"dark-text-input no-box-shadow",as:"select",custom:!0,onChange:function(t){e.setState({target_video_codec:t.target.value})},value:this.state.target_video_codec,children:[Object(k.jsx)("option",{value:"HEVC",children:"H.265 (HEVC)"}),Object(k.jsx)("option",{value:"AVC",children:"H.264 (AVC)"}),Object(k.jsx)("option",{value:"VP9",children:"VP9"})]})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Create Stereo Audio Track"})}),Object(k.jsx)(T.a.Checkbox,{"aria-label":"Create Stereo Audio Track Checkbox",onChange:function(t){e.setState({create_stereo_audio:t.target.checked})},checked:this.state.create_stereo_audio})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Skip HDR"})}),Object(k.jsx)(T.a.Checkbox,{"aria-label":"Skip HDR Checkbox",onChange:function(t){e.setState({skip_hdr:t.target.checked})},checked:this.state.skip_hdr})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Path Masks"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"Plex Versions,private,.m4a","aria-label":"path_masks","aria-describedby":"basic-addon1",onChange:function(t){e.setState({path_masks:t.target.value})},value:this.state.path_masks})]})]}),Object(k.jsxs)(v.a.Footer,{children:[Object(k.jsx)(p.a,{variant:"secondary",onClick:this.props.closeHandler,children:"Close"}),Object(k.jsx)(p.a,{variant:"primary",onClick:this.submitLib,children:"Create"})]})]})})}}]),a}(i.a.Component),L=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={folder:e.folder,priority:e.priority,fs_check_interval:e.fs_check_interval,path_masks:e.path_masks,target_video_codec:e.target_video_codec,create_stereo_audio:e.create_stereo_audio,skip_hdr:e.skip_hdr},s.putChanges=s.putChanges.bind(Object(N.a)(s)),s.deleteLibrary=s.deleteLibrary.bind(Object(N.a)(s)),s}return Object(o.a)(a,[{key:"putChanges",value:function(){var e=this,t={folder:this.state.folder,priority:parseInt(this.state.priority),fs_check_interval:this.state.fs_check_interval,path_masks:this.state.path_masks.split(","),command_decider_settings:JSON.stringify({target_video_codec:this.state.target_video_codec,create_stereo_audio:this.state.create_stereo_audio,skip_hdr:this.state.skip_hdr})};u.a.put("/api/web/v1/library/".concat(this.props.id),t).then((function(){e.props.closeHandler()})).catch((function(t){console.error("/api/web/v1/library/".concat(e.props.id," failed with error: ").concat(t))}))}},{key:"deleteLibrary",value:function(){var e=this;u.a.delete("/api/web/v1/library/".concat(this.props.id)).then((function(){e.props.closeHandler()})).catch((function(t){console.error("/api/web/v1/library/".concat(e.props.id," failed with error: ").concat(t))}))}},{key:"render",value:function(){var e=this;return Object(k.jsx)("div",{children:Object(k.jsxs)(v.a,{show:this.props.show,onHide:this.props.closeHandler,children:[Object(k.jsx)(v.a.Header,{closeButton:!0,children:Object(k.jsx)(v.a.Title,{children:"Edit Library"})}),Object(k.jsxs)(v.a.Body,{children:[Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Folder"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"/home/user/lib1","aria-label":"folder","aria-describedby":"basic-addon1",onChange:function(t){e.setState({folder:t.target.value})},value:this.state.folder})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Priority"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"0","aria-label":"priority","aria-describedby":"basic-addon1",onChange:function(t){e.setState({priority:t.target.value})},value:this.state.priority})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"File System Check Interval"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"fs_check_interval","aria-describedby":"basic-addon1",onChange:function(t){e.setState({fs_check_interval:t.target.value})},value:this.state.fs_check_interval})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Target Video Codec"})}),Object(k.jsxs)(I.a,{className:"dark-text-input no-box-shadow",as:"select",custom:!0,onChange:function(t){e.setState({target_video_codec:t.target.value})},value:this.state.target_video_codec,children:[Object(k.jsx)("option",{value:"HEVC",children:"H.265 (HEVC)"}),Object(k.jsx)("option",{value:"AVC",children:"H.264 (AVC)"}),Object(k.jsx)("option",{value:"VP9",children:"VP9"})]})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Create Stereo Audio Track"})}),Object(k.jsx)(T.a.Checkbox,{"aria-label":"Create Stereo Audio Track Checkbox",onChange:function(t){e.setState({create_stereo_audio:t.target.checked})},checked:this.state.create_stereo_audio})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Skip HDR"})}),Object(k.jsx)(T.a.Checkbox,{"aria-label":"Skip HDR Checkbox",onChange:function(t){e.setState({skip_hdr:t.target.checked})},checked:this.state.skip_hdr})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Path Masks"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"Plex Versions,private,.m4a","aria-label":"path_masks","aria-describedby":"basic-addon1",onChange:function(t){e.setState({path_masks:t.target.value})},value:this.state.path_masks})]})]}),Object(k.jsxs)(v.a.Footer,{children:[Object(k.jsx)(p.a,{className:"delete-button",variant:"danger",onClick:this.deleteLibrary,children:"Delete"}),Object(k.jsx)(p.a,{variant:"secondary",onClick:this.props.closeHandler,children:"Close"}),Object(k.jsx)(p.a,{variant:"primary",onClick:this.putChanges,children:"Update"})]})]})})}}]),a}(i.a.Component),D=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(){return Object(c.a)(this,a),t.apply(this,arguments)}return Object(o.a)(a,[{key:"render",value:function(){var e=this.props.queue;null===e&&(e=[]);var t=e.map((function(e,t){return Object(k.jsx)(F,{index:t+1,path:e.path,command:e.command.join(" ")},e.uuid)}));return Object(k.jsx)("div",{children:Object(k.jsxs)(v.a,{show:this.props.show,onHide:this.props.closeHandler,size:"lg",children:[Object(k.jsx)(v.a.Header,{closeButton:!0,children:Object(k.jsx)(v.a.Title,{children:"Queue"})}),Object(k.jsx)(v.a.Body,{children:Object(k.jsxs)(H.a,{children:[Object(k.jsx)("thead",{children:Object(k.jsxs)("tr",{children:[Object(k.jsx)("th",{scope:"col",children:"#"}),Object(k.jsx)("th",{scope:"col",children:"File"}),Object(k.jsx)("th",{scope:"col",children:"Cmd"})]})}),Object(k.jsx)("tbody",{children:t})]})}),Object(k.jsx)(v.a.Footer,{children:Object(k.jsx)(p.a,{variant:"secondary",onClick:this.props.closeHandler,children:"Close"})})]})})}}]),a}(i.a.Component);function F(e){return Object(k.jsxs)("tr",{children:[Object(k.jsx)("th",{scope:"row",children:e.index}),Object(k.jsx)("td",{children:e.path}),Object(k.jsx)("td",{children:Object(k.jsx)(C,{title:e.command})})]})}var M=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).timerID=void 0,s.state={jobs:[],waitingOnServer:!0},s.timerID=setTimeout((function(){}),Number.POSITIVE_INFINITY),clearInterval(s.timerID),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;this.tick(),this.timerID=setInterval((function(){return e.tick()}),2e3)}},{key:"componentWillUnmount",value:function(){clearInterval(this.timerID)}},{key:"tick",value:function(){var e=this;u.a.get("/api/web/v1/history").then((function(t){var a=t.data.history;void 0!==a?(a.reverse(),e.setState({jobs:a,waitingOnServer:!1})):console.error("Response from /api/web/v1/history returned undefined for data.history")})).catch((function(e){console.error("Request to /api/web/v1/history failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this.state.jobs.map((function(e,t){return Object(k.jsx)(A,{datetime:e.datetime_completed,file:e.file},t)})),t=Object(k.jsxs)("tr",{children:[Object(k.jsx)("th",{scope:"row",children:"-"}),Object(k.jsx)("td",{children:"Waiting on server"})]});return Object(k.jsx)(x.a,{children:Object(k.jsxs)(H.a,{hover:!0,size:"sm",children:[Object(k.jsx)("thead",{children:Object(k.jsxs)("tr",{children:[Object(k.jsx)("th",{scope:"col",children:"Time Completed"}),Object(k.jsx)("th",{scope:"col",children:"File"})]})}),Object(k.jsx)("tbody",{children:this.state.waitingOnServer?t:e})]})})}}]),a}(i.a.Component);function A(e){return Object(k.jsxs)("tr",{children:[Object(k.jsx)("td",{children:e.datetime}),Object(k.jsx)("td",{children:e.file})]})}a(88);var q=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={controller_version:"Could not contact a ring",web_api_versions:[],runner_api_versions:[]},s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;u.a.get("/api").then((function(t){e.setState({web_api_versions:t.data.web.versions,runner_api_versions:t.data.runner.versions})})).catch((function(e){console.error("Request to /api failed with error: ".concat(e))})),u.a.get("/version").then((function(t){e.setState({controller_version:t.data})})).catch((function(e){console.error("Request to /api failed with error: ".concat(e))}))}},{key:"render",value:function(){return Object(k.jsxs)(k.Fragment,{children:[Object(k.jsx)("h5",{children:"About Encodarr"}),Object(k.jsxs)("p",{children:[Object(k.jsx)("b",{children:"License:"})," This project is licensed under the Mozilla Public License 2.0 a copy of which can be found ",Object(k.jsx)("a",{href:"https://github.com/BrenekH/encodarr/blob/master/LICENSE",target:"_blank",rel:"noreferrer",children:"here"})]}),Object(k.jsxs)("p",{children:[Object(k.jsx)("b",{children:"Controller Version:"})," ",this.state.controller_version]}),Object(k.jsx)("p",{className:"list-title",children:Object(k.jsx)("b",{children:"Supported API Versions:"})}),Object(k.jsxs)("ul",{className:"api-list",children:[Object(k.jsxs)("li",{children:[Object(k.jsx)("b",{children:"Web:"})," ",this.state.web_api_versions.join(", ")]}),Object(k.jsxs)("li",{children:[Object(k.jsx)("b",{children:"Runner:"})," ",this.state.runner_api_versions.join(", ")]})]}),Object(k.jsxs)("p",{children:[Object(k.jsx)("b",{children:"GitHub Repository:"})," ",Object(k.jsx)("a",{href:"https://github.com/BrenekH/encodarr",target:"_blank",rel:"noreferrer",children:"https://github.com/BrenekH/encodarr"})]})]})}}]),a}(i.a.Component),B=(a(89),function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={inputValues:{healthCheckInterval:"",unresponsiveRunnerTimeout:"",logVerbosity:""},showSavedIndicator:!1},s.handleClick=s.handleClick.bind(Object(N.a)(s)),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){this.updateSettings()}},{key:"createChangeHandler",value:function(e){var t=this,a=arguments.length>1&&void 0!==arguments[1]&&arguments[1],s=function(s){var i=Object.assign({},t.state.inputValues);i[e]=a?s.target.checked:s.target.value,t.setState({inputValues:i})};return s.bind(this),s}},{key:"handleClick",value:function(){var e=this;u.a.put("/api/web/v1/settings",{HealthCheckInterval:this.state.inputValues.healthCheckInterval,HealthCheckTimeout:this.state.inputValues.unresponsiveRunnerTimeout,LogVerbosity:this.state.inputValues.logVerbosity}).then((function(t){t.status>=200&&t.status<=299?e.setState({showSavedIndicator:!0}):console.error(t),e.updateSettings()}))}},{key:"updateSettings",value:function(){var e=this;u.a.get("/api/web/v1/settings").then((function(t){e.setState({inputValues:{healthCheckInterval:t.data.HealthCheckInterval,unresponsiveRunnerTimeout:t.data.HealthCheckTimeout,logVerbosity:t.data.LogVerbosity}})}))}},{key:"render",value:function(){var e=this,t=this.state.showSavedIndicator?Object(k.jsx)(K,{}):null;return this.state.showSavedIndicator&&setTimeout((function(){e.setState({showSavedIndicator:!1})}),5e3),Object(k.jsxs)(k.Fragment,{children:[Object(k.jsxs)("div",{children:[Object(k.jsx)("h5",{children:"Runner Health"}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Runner Health Check Interval"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"health-check-interval","aria-describedby":"basic-addon1",onChange:this.createChangeHandler("healthCheckInterval"),value:this.state.inputValues.healthCheckInterval})]}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Unresponsive Runner Timeout"})}),Object(k.jsx)(I.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"unresponsive-runner-timeout","aria-describedby":"basic-addon1",onChange:this.createChangeHandler("unresponsiveRunnerTimeout"),value:this.state.inputValues.unresponsiveRunnerTimeout})]}),Object(k.jsx)("div",{className:"spacer"}),Object(k.jsx)("h5",{children:"Logging"}),Object(k.jsxs)(T.a,{className:"mb-3",children:[Object(k.jsx)(T.a.Prepend,{children:Object(k.jsx)(T.a.Text,{children:"Log Verbosity"})}),Object(k.jsxs)(I.a,{className:"dark-text-input no-box-shadow",as:"select",custom:!0,onChange:this.createChangeHandler("logVerbosity"),value:this.state.inputValues.logVerbosity,children:[Object(k.jsx)("option",{value:"TRACE",children:"Trace"}),Object(k.jsx)("option",{value:"DEBUG",children:"Debug"}),Object(k.jsx)("option",{value:"INFO",children:"Info"}),Object(k.jsx)("option",{value:"WARNING",children:"Warning"}),Object(k.jsx)("option",{value:"ERROR",children:"Error"}),Object(k.jsx)("option",{value:"CRITICAL",children:"Critical"})]})]}),Object(k.jsx)("div",{className:"smol-spacer"}),Object(k.jsx)(p.a,{variant:"light",onClick:this.handleClick,children:"Save"}),t]}),Object(k.jsx)("div",{className:"spacer"}),Object(k.jsx)(q,{})]})}}]),a}(i.a.Component));function K(){return Object(k.jsx)("p",{className:"pop-in-out",style:{display:"inline"},children:"Saved!"})}a(90);var W=a.p+"static/media/Encodarr-Logo.4b0cc1bf.svg";function Q(){return Object(k.jsxs)("div",{className:"header-flex header-content text-center",children:[Object(k.jsx)("img",{src:W,alt:"",height:"60px",title:""}),Object(k.jsx)("h1",{children:"ncodarr"})]})}var U=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(){return Object(c.a)(this,a),t.apply(this,arguments)}return Object(o.a)(a,[{key:"handleSelect",value:function(e){switch(e){case"libraries":window.history.replaceState(void 0,"","/libraries"),document.title="Libraries - Encodarr";break;case"history":window.history.replaceState(void 0,"","/history"),document.title="History - Encodarr";break;case"settings":window.history.replaceState(void 0,"","/settings"),document.title="Settings - Encodarr";break;case"running":window.history.replaceState(void 0,"","/running"),document.title="Running - Encodarr"}}},{key:"render",value:function(){var e="running";switch(window.location.pathname){case"/libraries":e="libraries";break;case"/history":e="history";break;case"/settings":e="settings"}return Object(k.jsxs)("div",{className:"container",children:[Object(k.jsx)(Q,{}),Object(k.jsxs)(h.a.Container,{id:"tab-nav",defaultActiveKey:e,transition:!1,onSelect:this.handleSelect,children:[Object(k.jsxs)(j.a,{fill:!0,variant:"pills",children:[Object(k.jsx)(j.a.Item,{children:Object(k.jsx)(j.a.Link,{eventKey:"running",children:"Running"})}),Object(k.jsx)(j.a.Item,{children:Object(k.jsx)(j.a.Link,{eventKey:"libraries",children:"Libraries"})}),Object(k.jsx)(j.a.Item,{children:Object(k.jsx)(j.a.Link,{eventKey:"history",children:"History"})}),Object(k.jsx)(j.a.Item,{children:Object(k.jsx)(j.a.Link,{eventKey:"settings",children:"Settings"})})]}),Object(k.jsx)("div",{className:"spacer"}),Object(k.jsxs)(h.a.Content,{children:[Object(k.jsx)(h.a.Pane,{eventKey:"running",mountOnEnter:!0,unmountOnExit:!0,children:Object(k.jsx)(w,{})}),Object(k.jsx)(h.a.Pane,{eventKey:"libraries",mountOnEnter:!0,unmountOnExit:!0,children:Object(k.jsx)(R,{})}),Object(k.jsx)(h.a.Pane,{eventKey:"history",mountOnEnter:!0,unmountOnExit:!0,children:Object(k.jsx)(M,{})}),Object(k.jsx)(h.a.Pane,{eventKey:"settings",mountOnEnter:!0,unmountOnExit:!0,children:Object(k.jsx)(B,{})})]})]}),Object(k.jsx)("div",{className:"smol-spacer"})]})}}]),a}(i.a.Component),J=function(e){e&&e instanceof Function&&a.e(3).then(a.bind(null,93)).then((function(t){var a=t.getCLS,s=t.getFID,i=t.getFCP,n=t.getLCP,r=t.getTTFB;a(e),s(e),i(e),n(e),r(e)}))};r.a.render(Object(k.jsx)(i.a.StrictMode,{children:Object(k.jsx)(U,{})}),document.getElementById("root")),J()}},[[91,1,2]]]);
//# sourceMappingURL=main.0078a14a.chunk.js.map