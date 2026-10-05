import { SafeAreaView, } from "react-native-safe-area-context";
import { Text, View, TouchableOpacity, TextInput, Modal, } from "react-native";
import { useNavigation } from "@react-navigation/native";
import Styled from "./style";
import { useState } from "react";
import { Feather } from "react-native-feather";




export default () => {

    const [imagePassowrd, setImagePassowrd] = useState("eye-off");
    const [passowrdFild, setPassowrdFild] = useState("123");
    const [emailFild, setEmailFild] = useState('');
    const [seePassowrd, setSeePassowrd] = useState(true);
    const [loading, setLoading] = useState(false);

    const Navigator = useNavigation();



    function sair() {
        Navigator.reset({
            routes: [{ name: 'SingIn' }]
        })
    }

    return (
        <SafeAreaView style={Styled.safeArea}>
            <View style={''}>
                <TextInput 
                    style={''}
                    placeholder="Digitel seu login/Email"
                    placeholderTextColor='#5E6737'
                    value={emailFild}
                    onChangeText={t => setEmailFild(t)}
                />
            </View>
            
            <View style={''}>
                <TextInput 
                    style={''}
                    placeholder="Digite sua Senha"
                    placeholderTextColor='#5E6737'
                    value={passowrdFild}
                    secureTextEntry={seePassowrd}
                    onChangeText={t => setPassowrdFild(t)}
                />
                <TouchableOpacity  style={{ marginRight: 12 }} >
                    <Feather name={imagePassowrd} size={20} color="yellow"/>
                </TouchableOpacity>

                <TouchableOpacity style={""} >
                    <Text style={""}>Login</Text>
                </TouchableOpacity>

                <TouchableOpacity onPress={''} style={''} >
                    <Text>Ainda não possui uma conta?</Text>
                    <Text style={''}>Cadastre-se</Text>
                </TouchableOpacity>
            </View>

            <Modal
                transparent={true}
                animationType="none"
                visible={loading}
            />

        </SafeAreaView>
    );
};