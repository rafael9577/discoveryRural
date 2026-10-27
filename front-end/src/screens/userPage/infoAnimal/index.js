import React from "react";
import { View, Text, TouchableOpacity, Image } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { Ionicons } from "@expo/vector-icons";

import styled from "./styled";

export default ({ navigation }) => {
  return (
    <SafeAreaView style={styled.safeArea}>
      {/* Cabeçalho */}
      <View style={styled.containerCabecalho}>
        <TouchableOpacity onPress={() => navigation?.goBack()}>
          <Ionicons name="arrow-back" size={28} color="#FFFFFF" />
        </TouchableOpacity>
        <Text style={styled.tituloCabecalho}>Informações do animal</Text>
      </View>

      {/* Cartão Principal */}
      <View style={styled.suporteCartao}>
        <View style={styled.containerPrincipal}>
          
          {/* Bloco Superior */}
          <View style={styled.blocoSuperior}>
            <Image
              source={{
                uri: "https://www.infoescola.com/wp-content/uploads/2018/03/sagui-de-tufos-brancos_790108426-1000x667.jpg",
              }}
              style={styled.imagemAnimal}
            />

            <View style={styled.divisorVertical} />

            <View style={styled.containerInformacoes}>
              <Text style={styled.nomeAnimal}>SAGUI-DE-TUFOS-BRANCOS</Text>
              <View style={styled.divisorHorizontal} />
              <Text style={styled.nomeCientifico}>(Callithrix jacchus)</Text>
              <Text style={styled.textoInformacao}>Peso médio:</Text>
              <Text style={styled.textoInformacao}>Tamanho médio:</Text>
            </View>
          </View>

        </View>
      </View>
    </SafeAreaView>
  );
};