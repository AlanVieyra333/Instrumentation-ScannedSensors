//var s1-1-p = false;

function updateSensorData(){
	//alert("UP");
	var reqJSON = {
		Operation:  1
	};	
	//reqJSON.Operation = 1;
	
	var data = JSON.stringify(reqJSON);
	var ok = function(res){
		var degrees = res.Degrees;
		var riego = res.Motor;
		var bpm = res.BPM;
		
		if(res.Code == 1){	// Datos recividos correctamente.
			$('#s1-info').html(degrees);	/*
			if(degrees  >= 10 && degrees <= 25){
				$('#img-s1-1').removeAttr('src');
				$('#img-s1-1').attr('src','img/animations/ventilador.gif');
				$('#img-s1-1').show();
				
				$('#img-s1-2').removeAttr('src');
				$('#img-s1-2').attr('src','img/animations/ventilador.png');
				$('#img-s1-2').show();
			}else if(degrees  >= 35 && degrees <= 60){
				$('#img-s1-1').removeAttr('src');
				$('#img-s1-1').attr('src','img/animations/ventilador.png');
				$('#img-s1-1').show();
				
				$('#img-s1-2').removeAttr('src');
				$('#img-s1-2').attr('src','img/animations/ventilador.gif');
				$('#img-s1-2').show();
			}else{
				$('#img-s1-1').removeAttr('src');
				$('#img-s1-1').attr('src','img/animations/ventilador.png');
				$('#img-s1-1').show();
				
				$('#img-s1-2').removeAttr('src');
				$('#img-s1-2').attr('src','img/animations/ventilador.png');
				$('#img-s1-2').show();
			}//*/	
			
			$('#s2-info').html(riego);	/*
			if(riego <= 30){
				$('#img-s2-1').removeAttr('src');
				$('#img-s2-1').attr('src','img/animations/riego.gif');
				$('#img-s2-1').show();
			}else{
				$('#img-s2-1').removeAttr('src');
				$('#img-s2-1').attr('src','img/animations/riego.png');
				$('#img-s2-1').show();
			}//*/
		
			$('#s3-info').html(bpm);	/*
			if(bpm > 0){
				$('#img-s3-1').removeAttr('src');
				$('#img-s3-1').attr('src','img/animations/cardio.gif');
				$('#img-s3-1').show();
			}else{
				$('#img-s3-1').removeAttr('src');
				$('#img-s3-1').attr('src','img/animations/cardio.png');
				$('#img-s3-1').show();
			}//*/
			
		}else{
			alert(res.Message);
		}		
	};
	
	sendAjaxJSON(data, ok, false);
}